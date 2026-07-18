package kb

import "time"

// Document represents a single knowledge base document.
type Document struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Path       string    `json:"path"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
	ModifiedAt time.Time `json:"modified_at"`
	ChunkCount int       `json:"chunk_count"`
}

// Chunk represents a portion of a document.
type Chunk struct {
	ID         string    `json:"id"`
	DocumentID string    `json:"document_id"`
	Index      int       `json:"index"`
	Content    string    `json:"content"`
	Embedding  []float32 `json:"embedding,omitempty"`
}
