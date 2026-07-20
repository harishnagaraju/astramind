package ai

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOpenAIProviderInvalidJSON(t *testing.T) {

	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {

				w.Header().Set(
					"Content-Type",
					"application/json",
				)

				_, _ = w.Write([]byte("{ invalid json"))
			},
		),
	)

	defer server.Close()

	provider := &OpenAIProvider{
		baseURL: server.URL,
	}

	_, err := provider.Chat(
		ChatRequest{
			Model:  "dummy",
			APIKey: "dummy",
		},
	)

	if err == nil {
		t.Fatal("expected JSON decoding error")
	}
}
