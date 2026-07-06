package ai

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOpenAIProviderHTTPErrors(t *testing.T) {

	tests := []struct {
		name string
		code int
		body string
	}{
		{
			name: "BadRequest",
			code: http.StatusBadRequest,
			body: `{"error":{"message":"bad request"}}`,
		},
		{
			name: "Unauthorized",
			code: http.StatusUnauthorized,
			body: `{"error":{"message":"invalid api key"}}`,
		},
		{
			name: "Forbidden",
			code: http.StatusForbidden,
			body: `{"error":{"message":"forbidden"}}`,
		},
		{
			name: "TooManyRequests",
			code: http.StatusTooManyRequests,
			body: `{"error":{"message":"quota exceeded"}}`,
		},
		{
			name: "InternalServerError",
			code: http.StatusInternalServerError,
			body: `{"error":{"message":"internal error"}}`,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			server := httptest.NewServer(
				http.HandlerFunc(
					func(w http.ResponseWriter, r *http.Request) {

						w.WriteHeader(tt.code)
						_, _ = w.Write([]byte(tt.body))
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
				t.Fatalf(
					"expected error for HTTP %d",
					tt.code,
				)
			}
		})
	}
}
