package ai

import "strings"

// NewProvider creates a provider instance based on the configuration.
func NewProvider(cfg ProviderConfig) Provider {

	switch strings.ToLower(
		strings.TrimSpace(cfg.Provider),
	) {

	case "mock":
		return &MockProvider{}

	case "openai":
		return &OpenAIProvider{
			baseURL: cfg.BaseURL,
		}

	case "ollama":
		return &OllamaProvider{
			baseURL: "http://localhost:11434",
		}
	}

	if strings.TrimSpace(cfg.APIKey) == "" {
		return &MockProvider{}
	}

	return &OpenAIProvider{
		baseURL: cfg.BaseURL,
	}
}
