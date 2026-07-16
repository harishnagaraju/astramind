package ai

import (
	"context"
	"encoding/json"
	"github.com/harishnagaraju/astramind/internal/infrastructure/models"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOllamaProviderStreamIntegration(t *testing.T) {

	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if r.Method != http.MethodPost {
				t.Fatalf("expected POST, got %s", r.Method)
			}

			if r.URL.Path != "/api/chat" {
				t.Fatalf("expected /api/chat, got %s", r.URL.Path)
			}

			var req OllamaChatRequest

			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				t.Fatalf("decode request: %v", err)
			}

			if !req.Stream {
				t.Fatal("expected stream=true")
			}

			w.Header().Set(
				"Content-Type",
				"application/json",
			)

			encoder := json.NewEncoder(w)

			_ = encoder.Encode(OllamaStreamResponse{
				Message: OllamaStreamMessage{
					Content: "Hello ",
				},
				Done: false,
			})

			_ = encoder.Encode(OllamaStreamResponse{
				Message: OllamaStreamMessage{
					Content: "World",
				},
				Done: false,
			})

			_ = encoder.Encode(OllamaStreamResponse{
				Done: true,
			})
		}),
	)

	defer server.Close()

	provider := &OllamaProvider{
		baseURL: server.URL,
		model:   "gemma3:1b",
	}

	stream, err := provider.Stream(
		context.Background(),
		ChatRequest{
			Messages: []models.Message{
				{
					Role:    "user",
					Content: "hello",
				},
			},
		},
	)

	if err != nil {
		t.Fatalf("Stream() returned error: %v", err)
	}

	var tokens []string
	done := false

	for event := range stream.Events() {

		switch event.Type {

		case StreamEventToken:
			tokens = append(tokens, event.Content)

		case StreamEventDone:
			done = true

		case StreamEventError:
			t.Fatalf("unexpected stream error: %v", event.Err)
		}
	}

	if len(tokens) != 2 {
		t.Fatalf(
			"expected 2 tokens, got %d",
			len(tokens),
		)
	}

	if tokens[0] != "Hello " {
		t.Fatalf(
			"unexpected first token: %q",
			tokens[0],
		)
	}

	if tokens[1] != "World" {
		t.Fatalf(
			"unexpected second token: %q",
			tokens[1],
		)
	}

	if !done {
		t.Fatal("expected StreamEventDone")
	}
}
