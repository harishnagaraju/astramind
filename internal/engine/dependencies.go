package engine

import (
	"github.com/harishnagaraju/astramind/internal/features/chat"
	"github.com/harishnagaraju/astramind/internal/features/export"
	"github.com/harishnagaraju/astramind/internal/features/history"
	"github.com/harishnagaraju/astramind/internal/features/kb"
	"github.com/harishnagaraju/astramind/internal/features/search"
	"github.com/harishnagaraju/astramind/internal/features/session"
	"github.com/harishnagaraju/astramind/internal/infrastructure/ai"
)

type Dependencies struct {
	ProviderManager *ai.ProviderManager

	KnowledgeBase *kb.Manager

	ChatService *chat.Service

	HistoryService *history.Service

	SessionService *session.Service

	ExportService *export.Service

	SearchService *search.Service
}
