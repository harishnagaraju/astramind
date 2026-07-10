package ai

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOllamaProviderHTTPErrors(t *testing.T) {

	tests := []struct {
		name       string
		statusCode int
	}{
		{
			name:       "BadRequest",
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "Unauthorized",
			statusCode: http.StatusUnauthorized,
		},
		{
			name:       "NotFound",
			statusCode: http.StatusNotFound,
		},
		{
			name:       "TooManyRequests",
			statusCode: http.StatusTooManyRequests,
		},
		{
			name:       "InternalServerError",
			statusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {

		t.Run(
			tt.name,
			func(t *testing.T) {

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

							w.WriteHeader(
								tt.statusCode,
							)

							fmt.Fprint(
								w,
								`{"error":"test error"}`,
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
					t.Fatalf(
						"expected error for HTTP %d",
						tt.statusCode,
					)
				}
			},
		)
	}
}
