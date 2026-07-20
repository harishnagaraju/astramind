package ai

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/harishnagaraju/astramind/internal/infrastructure/models"
)

// TestOllamaRequestSetsContextWindow guards against a real bug found
// in manual testing: buildOllamaRequest previously sent no options
// at all, so Ollama fell back to its own default context window
// (commonly 2048 tokens). RAG prompts (system instructions +
// retrieved chunks + question + answer) routinely exceed that,
// causing answers to truncate mid-generation rather than error.
func TestOllamaRequestSetsContextWindow(t *testing.T) {

	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			var req OllamaChatRequest

			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				t.Fatalf("decode request: %v", err)
			}

			if req.Options == nil {
				t.Fatal("expected options to be set, got nil")
			}

			if req.Options.NumCtx < 4096 {
				t.Fatalf(
					"expected a context window of at least 4096, got %d",
					req.Options.NumCtx,
				)
			}

			w.Header().Set("Content-Type", "application/json")

			json.NewEncoder(w).Encode(OllamaChatResponse{
				Message: OllamaChatMessage{
					Role:    "assistant",
					Content: "ok",
				},
				Done: true,
			})
		}),
	)
	defer server.Close()

	provider := &OllamaProvider{
		baseURL: server.URL,
		model:   "gemma3:1b",
	}

	_, err := provider.Chat(ChatRequest{
		Messages: []models.Message{
			{Role: "user", Content: "hello"},
		},
	})

	if err != nil {
		t.Fatal(err)
	}
}
