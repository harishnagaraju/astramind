package ai

import "strings"

func NewProvider(
	cfg ProviderConfig,
) Provider {

	switch strings.ToLower(
		strings.TrimSpace(cfg.Provider),
	) {

	case "mock":
		return &MockProvider{}

	case "openai":
		return &OpenAIProvider{
			baseURL: cfg.BaseURL,
		}
	}

	if strings.TrimSpace(cfg.APIKey) == "" {
		return &MockProvider{}
	}

	return &OpenAIProvider{
		baseURL: cfg.BaseURL,
	}
}
