package ai

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestOpenAIProvider_Embed mirrors TestOllamaProvider_Embed's approach
// (httptest.Server, no injected client field needed - OpenAIProvider
// has only a baseURL field, confirmed against the real struct) but
// exercises OpenAI's genuinely different response shape: a nested
// "data": [{"embedding": [...]}] array, not Ollama's flat
// "embedding": [...]. Using the same test structure for a
// structurally different API would have missed real bugs specific to
// this parsing path.
func TestOpenAIProvider_Embed(t *testing.T) {
	tests := []struct {
		name           string
		inputText      string
		apiKey         string
		mockResponse   string
		mockStatusCode int
		expectError    bool
		expectedVector []float32
	}{
		{
			name:           "Valid embedding response",
			inputText:      "AstraMind Core Engine",
			apiKey:         "sk-test-key",
			mockResponse:   `{"data": [{"embedding": [0.123, -0.456, 0.789]}]}`,
			mockStatusCode: http.StatusOK,
			expectError:    false,
			expectedVector: []float32{0.123, -0.456, 0.789},
		},
		{
			name:           "Empty data array is treated as an error",
			inputText:      "",
			apiKey:         "sk-test-key",
			mockResponse:   `{"data": []}`,
			mockStatusCode: http.StatusOK,
			// Confirmed against the real code: len(result.Data) == 0
			// explicitly returns an error ("openai returned no
			// embedding data"), same pattern as the Ollama provider's
			// empty-embedding case.
			expectError:    true,
			expectedVector: nil,
		},
		{
			name:           "OpenAI backend returns HTTP 401 (bad API key)",
			inputText:      "Fail Processing",
			apiKey:         "invalid-key",
			mockResponse:   `{"error": {"message": "Invalid API key"}}`,
			mockStatusCode: http.StatusUnauthorized,
			expectError:    true,
			expectedVector: nil,
		},
		{
			name:           "Malformed JSON response body",
			inputText:      "Malformed Payload",
			apiKey:         "sk-test-key",
			mockResponse:   `{"data": [{"embedding": [0.1, unparsed_float_error]}]}`,
			mockStatusCode: http.StatusOK,
			expectError:    true,
			expectedVector: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Header.Get("Content-Type") != "application/json" {
					t.Errorf("expected application/json content type, got: %s", r.Header.Get("Content-Type"))
				}
				expectedAuth := "Bearer " + tt.apiKey
				if r.Header.Get("Authorization") != expectedAuth {
					t.Errorf("expected Authorization header %q, got: %q", expectedAuth, r.Header.Get("Authorization"))
				}
				if r.URL.Path != "/embeddings" {
					t.Errorf("expected request to /embeddings, got: %s", r.URL.Path)
				}
				w.WriteHeader(tt.mockStatusCode)
				fmt.Fprint(w, tt.mockResponse)
			}))
			defer server.Close()

			provider := &OpenAIProvider{baseURL: server.URL}

			vector, err := provider.Embed(EmbeddingRequest{Text: tt.inputText, APIKey: tt.apiKey})

			if (err != nil) != tt.expectError {
				t.Fatalf("unexpected error state: expectError=%v, got err=%v", tt.expectError, err)
			}

			if tt.expectError {
				return
			}

			if len(vector) != len(tt.expectedVector) {
				t.Fatalf("vector length mismatch: expected %d, got %d", len(tt.expectedVector), len(vector))
			}
			for i, val := range vector {
				if val != tt.expectedVector[i] {
					t.Errorf("vector[%d]: expected %f, got %f", i, tt.expectedVector[i], val)
				}
			}
		})
	}
}
