package ai

import (
	"context"
	"testing"
)

func TestMockProviderStream(t *testing.T) {

	provider := &MockProvider{}

	stream, err := provider.Stream(
		context.Background(),
		ChatRequest{},
	)

	if err != nil {
		t.Fatalf(
			"unexpected error: %v",
			err,
		)
	}

	var events []StreamEvent

	for event := range stream.Events() {
		events = append(events, event)
	}

	if len(events) != 5 {
		t.Fatalf(
			"expected 5 events, got %d",
			len(events),
		)
	}
	expectedTokens := []string{
		"Hello",
		" from",
		" Mock",
		" AI!",
	}

	for i, token := range expectedTokens {

		if events[i].Type != StreamEventToken {
			t.Fatalf(
				"event %d: expected %q, got %q",
				i,
				StreamEventToken,
				events[i].Type,
			)
		}

		if events[i].Content != token {
			t.Fatalf(
				"event %d: expected token %q, got %q",
				i,
				token,
				events[i].Content,
			)
		}
	}

	if events[4].Type != StreamEventDone {
		t.Fatalf(
			"expected final event %q, got %q",
			StreamEventDone,
			events[4].Type,
		)
	}

}
