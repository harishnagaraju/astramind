package kb

// Embedder generates a vector embedding for a piece of text. It is
// defined locally in the kb package (rather than importing the ai
// package directly) to keep kb decoupled from the AI provider layer,
// the same separation already used between search/history/storage.
type Embedder interface {
	Embed(text string) ([]float32, error)
}
