package ai

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOllamaProviderInvalidJSON(t *testing.T) {

	server := httptest.NewServer(
		http.HandlerFunc(
			func(
				w http.ResponseWriter,
				r *http.Request,
			) {

				w.Header().Set(
					"Content-Type",
					"application/json",
				)

				fmt.Fprint(
					w,
					`{ invalid json`,
				)
			},
		),
	)

	defer server.Close()

	provider := &OllamaProvider{
		baseURL: server.URL,
		model:   "gemma3:1b",
	}

	_, err := provider.Chat(
		ChatRequest{
			Model: "gemma3:1b",
		},
	)

	if err == nil {
		t.Fatal(
			"expected JSON decoding error",
		)
	}
}
