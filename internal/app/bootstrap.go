package app

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"

	"github.com/harishnagaraju/astramind/internal/ai"
	"github.com/harishnagaraju/astramind/internal/chat"
	"github.com/harishnagaraju/astramind/internal/kb"
)

func (a *App) initialize() error {

	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("could not load .env file: %w", err)
	}

	a.providerName = strings.TrimSpace(
		os.Getenv("AI_PROVIDER"),
	)

	a.apiKey = os.Getenv("OPENAI_API_KEY")
	a.model = os.Getenv("OPENAI_MODEL")
	a.baseURL = os.Getenv("OPENAI_BASE_URL")

	if a.providerName != "ollama" && a.apiKey == "" {
		return fmt.Errorf("no OpenAI API key configured")
	}

	if a.model == "" {
		a.model = "gpt-4o-mini"
	}

	if a.baseURL == "" {
		a.baseURL = "https://api.openai.com/v1"
	}

	if a.providerName == "" {
		a.providerName = "openai"
	}

	provider := ai.NewProvider(
		ai.ProviderConfig{
			Provider: a.providerName,
			APIKey:   a.apiKey,
			Model:    a.model,
			BaseURL:  a.baseURL,
		},
	)

	a.manager = ai.NewProviderManager(provider)

	storage := kb.NewJSONStorage("data")

	a.kbManager = kb.NewManager(storage)

	a.service = chat.NewService(
		chat.Dependencies{
			ProviderManager: a.manager,
			KnowledgeBase:   a.kbManager,
		},
	)

	return nil
}
