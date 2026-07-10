package ai

import "testing"

func TestOllamaProviderConnectionFailure(t *testing.T) {

	provider := &OllamaProvider{
		baseURL: "http://127.0.0.1:65535",
		model:   "gemma3:1b",
	}

	_, err := provider.Chat(
		ChatRequest{
			Model: "gemma3:1b",
		},
	)

	if err == nil {
		t.Fatal(
			"expected connection failure",
		)
	}
}
