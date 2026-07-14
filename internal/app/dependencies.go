package app

import (
	"github.com/harishnagaraju/astramind/internal/ai"
	"github.com/harishnagaraju/astramind/internal/chat"
	"github.com/harishnagaraju/astramind/internal/kb"
)

// Dependencies contains all long-lived application services.
type Dependencies struct {
	ProviderManager *ai.ProviderManager
	KnowledgeBase   *kb.Manager
	ChatService     *chat.Service
}
