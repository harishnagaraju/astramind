package ai

import (
	"context"
	"errors"
	"testing"
)

// nonStreamingProvider implements Provider but deliberately NOT
// StreamingProvider or EmbeddingProvider - used to test
// ProviderManager's behavior when the current provider lacks a
// capability entirely (not a runtime failure, a missing interface).
type nonStreamingProvider struct{}

func (p *nonStreamingProvider) Name() string { return "NonStreaming" }
func (p *nonStreamingProvider) Chat(request ChatRequest) (string, error) {
	return "chat reply", nil
}

// failingStreamProvider implements StreamingProvider but always
// returns an error from Stream() - used to test what happens on a
// genuine runtime failure (as opposed to a missing capability).
type failingStreamProvider struct{}

func (p *failingStreamProvider) Name() string { return "FailingStream" }
func (p *failingStreamProvider) Chat(request ChatRequest) (string, error) {
	return "", errors.New("chat not used in this test")
}
func (p *failingStreamProvider) Stream(ctx context.Context, request ChatRequest) (Stream, error) {
	return nil, errors.New("simulated stream failure")
}

// failingEmbedProvider implements EmbeddingProvider but always
// returns an error from Embed().
type failingEmbedProvider struct{}

func (p *failingEmbedProvider) Name() string { return "FailingEmbed" }
func (p *failingEmbedProvider) Chat(request ChatRequest) (string, error) {
	return "", errors.New("chat not used in this test")
}
func (p *failingEmbedProvider) Embed(request EmbeddingRequest) ([]float32, error) {
	return nil, errors.New("simulated embed failure")
}

// TestProviderManager_Stream_NoCapability confirms Stream() returns
// an immediate, clear error when the current provider doesn't
// implement StreamingProvider at all - no fallback attempt.
func TestProviderManager_Stream_NoCapability(t *testing.T) {
	pm := NewProviderManager(&nonStreamingProvider{})

	_, err := pm.Stream(context.Background(), ChatRequest{})

	if err == nil {
		t.Fatal("expected an error when the provider doesn't support streaming, got nil")
	}
}

// TestProviderManager_Stream_RuntimeFailureHasNoFallback is the real,
// verified finding this test suite exists to document: unlike Chat(),
// Stream() has NO fallback logic of its own. A genuine runtime
// failure from a provider that DOES implement StreamingProvider is
// returned directly - ProviderManager never attempts pm.fallback for
// Stream(), even though it does exactly that for Chat(). This is a
// real behavioral fact about the current code, confirmed by reading
// Stream()'s implementation (it only type-asserts and calls through -
// no error branch reaches pm.fallback at all), not an assumption.
func TestProviderManager_Stream_RuntimeFailureHasNoFallback(t *testing.T) {
	pm := NewProviderManager(&failingStreamProvider{})

	_, err := pm.Stream(context.Background(), ChatRequest{})

	if err == nil {
		t.Fatal("expected the simulated stream failure to propagate, got nil")
	}

	// Confirm no failover occurred - the primary provider is still
	// the one that failed, not pm.fallback (MockProvider).
	if pm.ProviderName() != "FailingStream" {
		t.Fatalf("expected provider to remain 'FailingStream' (no fallback for Stream), got %q", pm.ProviderName())
	}
}

// TestProviderManager_Embed_NoCapability mirrors the Stream test for
// Embed() - immediate error, no fallback, when the provider doesn't
// implement EmbeddingProvider.
func TestProviderManager_Embed_NoCapability(t *testing.T) {
	pm := NewProviderManager(&nonStreamingProvider{})

	_, err := pm.Embed(EmbeddingRequest{Text: "test"})

	if err == nil {
		t.Fatal("expected an error when the provider doesn't support embeddings, got nil")
	}
}

// TestProviderManager_Embed_RuntimeFailureHasNoFallback is Embed()'s
// equivalent of the Stream finding above: a genuine runtime failure
// from a provider that DOES implement EmbeddingProvider propagates
// directly, with no attempt to fall back to pm.fallback.
func TestProviderManager_Embed_RuntimeFailureHasNoFallback(t *testing.T) {
	pm := NewProviderManager(&failingEmbedProvider{})

	_, err := pm.Embed(EmbeddingRequest{Text: "test"})

	if err == nil {
		t.Fatal("expected the simulated embed failure to propagate, got nil")
	}

	if pm.ProviderName() != "FailingEmbed" {
		t.Fatalf("expected provider to remain 'FailingEmbed' (no fallback for Embed), got %q", pm.ProviderName())
	}
}

// TestProviderManager_Chat_FailoverIsPermanent documents Chat()'s
// actual, different behavior for contrast: on failure, it PERMANENTLY
// reassigns pm.provider to pm.fallback (confirmed: "pm.provider =
// pm.fallback" mutates state, not just this one call) - meaning a
// SUBSEQUENT Stream()/Embed() call after a Chat() failure WOULD then
// go through the fallback provider, just not because Stream/Embed
// triggered it themselves.
func TestProviderManager_Chat_FailoverIsPermanent(t *testing.T) {
	pm := NewProviderManager(&alwaysFailingChatProvider{})

	_, err := pm.Chat(ChatRequest{})
	if err != nil {
		t.Fatalf("expected Chat to succeed via fallback, got error: %v", err)
	}

	if pm.ProviderName() != "Mock AI" {
		t.Fatalf("expected provider to have permanently switched to the fallback after failure, got %q", pm.ProviderName())
	}

	// A second call, even with no new failure, stays on the fallback -
	// confirming the switch is permanent, not just for one retry.
	_, err = pm.Chat(ChatRequest{})
	if err != nil {
		t.Fatalf("expected second Chat call to also succeed via the now-permanent fallback, got: %v", err)
	}
	if pm.ProviderName() != "Mock AI" {
		t.Fatalf("expected provider to remain on fallback for subsequent calls, got %q", pm.ProviderName())
	}
}

type alwaysFailingChatProvider struct{}

func (p *alwaysFailingChatProvider) Name() string { return "AlwaysFailingChat" }
func (p *alwaysFailingChatProvider) Chat(request ChatRequest) (string, error) {
	return "", errors.New("simulated primary provider failure")
}
