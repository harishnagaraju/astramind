package ai

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOpenAIProviderUnauthorized(t *testing.T) {

	server := httptest.NewServer(
		http.HandlerFunc(
			func(
				w http.ResponseWriter,
				r *http.Request,
			) {

				w.WriteHeader(http.StatusUnauthorized)

				_, _ = w.Write([]byte(`{
					"error":{
						"message":"invalid api key"
					}
				}`))
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
			APIKey: "bad-key",
		},
	)

	if err == nil {
		t.Fatal("expected error")
	}
}
