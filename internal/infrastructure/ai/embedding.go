package ai

// EmbeddingRequest describes a request to embed a single piece of text.
type EmbeddingRequest struct {
	Model  string
	APIKey string
	Text   string
}

// EmbeddingProvider is implemented by providers that can generate
// vector embeddings for text. Not every Provider supports this, so
// callers type-assert against it - the same optional-capability
// pattern already used by StreamingProvider.
type EmbeddingProvider interface {
	Embed(
		request EmbeddingRequest,
	) ([]float32, error)
}
