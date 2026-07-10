package ai

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOllamaProviderChatIntegration(t *testing.T) {

	server := httptest.NewServer(
		http.HandlerFunc(
			func(
				w http.ResponseWriter,
				r *http.Request,
			) {

				if r.URL.Path != "/api/chat" {
					t.Fatalf(
						"unexpected path: %s",
						r.URL.Path,
					)
				}

				if r.Method != http.MethodPost {
					t.Fatalf(
						"unexpected method: %s",
						r.Method,
					)
				}

				w.Header().Set(
					"Content-Type",
					"application/json",
				)

				fmt.Fprint(
					w,
					`{
						"model":"gemma3:1b",
						"message":{
							"role":"assistant",
							"content":"Hello from Ollama test server"
						},
						"done":true
					}`,
				)
			},
		),
	)

	defer server.Close()

	provider := &OllamaProvider{
		baseURL: server.URL,
		model:   "gemma3:1b",
	}

	reply, err := provider.Chat(
		ChatRequest{
			Model: "gemma3:1b",
		},
	)

	if err != nil {
		t.Fatalf(
			"unexpected error: %v",
			err,
		)
	}

	if reply != "Hello from Ollama test server" {
		t.Fatalf(
			"unexpected reply: %q",
			reply,
		)
	}
}
