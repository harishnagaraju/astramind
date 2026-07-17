package engine

import (
	"github.com/harishnagaraju/astramind/internal/features/chat"
	"github.com/harishnagaraju/astramind/internal/features/kb"
	"github.com/harishnagaraju/astramind/internal/infrastructure/ai"
)

type Dependencies struct {
	ProviderManager *ai.ProviderManager

	KnowledgeBase *kb.Manager

	ChatService *chat.Service
}
