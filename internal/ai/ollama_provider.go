package ai

import "errors"

// OllamaProvider implements the Provider interface
// for locally hosted Ollama models.
type OllamaProvider struct {
	baseURL string
}

func (o *OllamaProvider) Name() string {
	return "Ollama"
}

func (o *OllamaProvider) Chat(
	request ChatRequest,
) (string, error) {
	return "", errors.New(
		"ollama provider not implemented",
	)
}
