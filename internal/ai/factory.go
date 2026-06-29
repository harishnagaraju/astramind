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
		return &OpenAIProvider{}
	}

	if strings.TrimSpace(cfg.APIKey) == "" {
		return &MockProvider{}
	}

	return &OpenAIProvider{}
}
