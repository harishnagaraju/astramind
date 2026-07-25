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

	// embeddingModel is passed through to every embedding request.
	// Previously this field didn't exist at all - EmbeddingRequest.Model
	// was always left empty, which silently made every embedding call
	// use ai.OpenAIProvider's hardcoded default ("text-embedding-3-small"),
	// completely independent of OPENAI_MODEL. That meant switching the
	// configured chat model (e.g. to an OpenRouter free model) had zero
	// effect on retrieval quality, with no visibility that a separate,
	// unconfigured, separately-billed embedding model was in use the
	// whole time. Now explicit and overridable via OPENAI_EMBEDDING_MODEL
	// (see bootstrap.go); left empty, behavior is unchanged from before -
	// this is purely additive, not a behavior change for anyone not
	// setting the new variable.
	embeddingModel string
}

// Embed implements kb.Embedder.
func (p *providerEmbedder) Embed(text string) ([]float32, error) {

	return p.providerManager.Embed(ai.EmbeddingRequest{
		APIKey: p.apiKey,
		Text:   text,
		Model:  p.embeddingModel,
	})
}
