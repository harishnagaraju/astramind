package app

import (
	"github.com/harishnagaraju/astramind/internal/ai"
	"github.com/harishnagaraju/astramind/internal/chat"
	"github.com/harishnagaraju/astramind/internal/kb"
)

type App struct {
	manager   *ai.ProviderManager
	kbManager *kb.Manager
	service   *chat.Service

	apiKey       string
	model        string
	baseURL      string
	providerName string
}

func New() *App {
	return &App{}
}
