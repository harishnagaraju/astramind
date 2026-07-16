package engine

import (
	"github.com/harishnagaraju/astramind/internal/features/chat"
	"github.com/harishnagaraju/astramind/internal/infrastructure/ai"
	"github.com/harishnagaraju/astramind/internal/kb"
)

// Dependencies contains all long-lived application services.
type Dependencies struct {
	ProviderManager *ai.ProviderManager
	KnowledgeBase   *kb.Manager
	ChatService     *chat.Service
}
