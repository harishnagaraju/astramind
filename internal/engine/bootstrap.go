package engine

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"

	"github.com/harishnagaraju/astramind/internal/features/chat"
	"github.com/harishnagaraju/astramind/internal/features/export"
	"github.com/harishnagaraju/astramind/internal/features/history"
	"github.com/harishnagaraju/astramind/internal/features/kb"
	"github.com/harishnagaraju/astramind/internal/features/search"
	"github.com/harishnagaraju/astramind/internal/features/session"
	"github.com/harishnagaraju/astramind/internal/infrastructure/ai"
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

	a.deps.ProviderManager = ai.NewProviderManager(provider)

	storage := kb.NewJSONStorage("data")

	a.deps.KnowledgeBase = kb.NewManager(storage)

	a.deps.ChatService = chat.NewService(
		chat.Dependencies{
			ProviderManager: a.deps.ProviderManager,
			KnowledgeBase:   a.deps.KnowledgeBase,
		},
	)

	a.deps.HistoryService = history.NewService()

	a.deps.SessionService = session.NewService()

	a.deps.ExportService = export.NewService()

	a.deps.SearchService = search.NewService(a.deps.HistoryService)

	a.dispatcher = newDispatcher(a)

	return nil
}
