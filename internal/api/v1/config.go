package v1

import (
	"github.com/harishnagaraju/astramind/internal/features/kb"
	"github.com/harishnagaraju/astramind/internal/infrastructure/ai"
)

type Config struct {
	ProviderName   string
	Model          string
	Version        string
	APIKey         string
	ProviderManager *ai.ProviderManager
	KnowledgeBase  *kb.Manager
}
