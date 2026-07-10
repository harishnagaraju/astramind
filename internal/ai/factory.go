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

		baseURL := strings.TrimSpace(cfg.BaseURL)
		if baseURL == "" {
			baseURL = "http://localhost:11434"
		}

		model := strings.TrimSpace(cfg.Model)
		if model == "" {
			model = "llama3"
		}

		return &OllamaProvider{
			baseURL: baseURL,
			model:   model,
		}
	}

	// Backward-compatible fallback behavior.
	if strings.TrimSpace(cfg.APIKey) == "" {
		return &MockProvider{}
	}

	return &OpenAIProvider{
		baseURL: cfg.BaseURL,
	}
}
