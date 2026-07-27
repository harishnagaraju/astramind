package ai

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestOllamaProvider_Embed uses httptest.Server rather than an
// injected http.RoundTripper - OllamaProvider has no client field to
// inject into (confirmed against the real struct: only baseURL and
// model), and Embed() constructs its own *http.Client{} internally
// on every call. httptest.Server needs none of that: it's a real,
// local, in-process HTTP server, and baseURL already is the one
// field Embed() uses to build its request - no production code
// change required at all.
func TestOllamaProvider_Embed(t *testing.T) {
	tests := []struct {
		name           string
		inputText      string
		mockResponse   string
		mockStatusCode int
		expectError    bool
		expectedVector []float32
	}{
		{
			name:           "Valid embedding response",
			inputText:      "AstraMind Core Engine",
			mockResponse:   `{"embedding": [0.123, -0.456, 0.789]}`,
			mockStatusCode: http.StatusOK,
			expectError:    false,
			expectedVector: []float32{0.123, -0.456, 0.789},
		},
		{
			name:           "Empty embedding array is treated as an error",
			inputText:      "",
			mockResponse:   `{"embedding": []}`,
			mockStatusCode: http.StatusOK,
			// Real behavior, confirmed by reading ollama_embedding.go:
			// an empty result.Embedding explicitly returns an error
			// ("ollama returned an empty embedding"), it does not
			// return an empty-but-successful vector. The proposal
			// this test replaces got this backwards (expected no
			// error for an empty embedding).
			expectError:    true,
			expectedVector: nil,
		},
		{
			name:           "Ollama backend returns HTTP 500",
			inputText:      "Fail Processing",
			mockResponse:   `Internal Server Error - Out of VRAM`,
			mockStatusCode: http.StatusInternalServerError,
			expectError:    true,
			expectedVector: nil,
		},
		{
			name:           "Malformed JSON response body",
			inputText:      "Malformed Payload",
			mockResponse:   `{"embedding": [0.123, -0.456, unparsed_float_error]}`,
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
				if r.URL.Path != "/api/embeddings" {
					t.Errorf("expected request to /api/embeddings, got: %s", r.URL.Path)
				}
				w.WriteHeader(tt.mockStatusCode)
				fmt.Fprint(w, tt.mockResponse) //nolint:errcheck // mock HTTP server response in a test, cannot meaningfully fail
			}))
			defer server.Close()

			provider := &OllamaProvider{baseURL: server.URL}

			vector, err := provider.Embed(EmbeddingRequest{Text: tt.inputText})

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