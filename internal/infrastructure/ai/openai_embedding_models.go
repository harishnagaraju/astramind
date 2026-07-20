package ai

// OpenAIEmbeddingRequest represents an embedding request sent to the
// OpenAI-compatible embeddings API.
type OpenAIEmbeddingRequest struct {
	Model string `json:"model"`
	Input string `json:"input"`
}

// OpenAIEmbeddingResponse represents the response returned by the
// OpenAI-compatible embeddings API.
type OpenAIEmbeddingResponse struct {
	Data []struct {
		Embedding []float32 `json:"embedding"`
	} `json:"data"`
}
