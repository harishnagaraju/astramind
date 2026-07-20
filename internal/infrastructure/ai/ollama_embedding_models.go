package ai

// OllamaEmbeddingRequest represents an embedding request sent to the
// Ollama API.
type OllamaEmbeddingRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

// OllamaEmbeddingResponse represents the response returned by the
// Ollama embeddings API.
type OllamaEmbeddingResponse struct {
	Embedding []float32 `json:"embedding"`
}
