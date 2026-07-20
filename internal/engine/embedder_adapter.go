package engine

import (
	"github.com/harishnagaraju/astramind/internal/infrastructure/ai"
)

// providerEmbedder adapts *ai.ProviderManager to kb.Embedder,
// translating a plain text string into the ai.EmbeddingRequest shape
// the provider layer expects. This is the one place allowed to know
// about both ai and kb - keeping that separation everywhere else in
// the codebase.
type providerEmbedder struct {
	providerManager *ai.ProviderManager
	apiKey          string
}

// Embed implements kb.Embedder.
func (p *providerEmbedder) Embed(text string) ([]float32, error) {

	return p.providerManager.Embed(ai.EmbeddingRequest{
		APIKey: p.apiKey,
		Text:   text,
	})
}
