package kb

// Repository provides higher-level knowledge base operations.
type Repository struct {
	manager *Manager
}

// NewRepository creates a new knowledge repository.
func NewRepository(manager *Manager) *Repository {
	return &Repository{
		manager: manager,
	}
}

// Documents returns all knowledge base documents.
func (r *Repository) Documents() ([]Document, error) {
	return r.manager.ListDocuments()
}

// Chunks returns all chunks for every stored document.
func (r *Repository) Chunks() ([]Chunk, error) {
	documents, err := r.manager.ListDocuments()
	if err != nil {
		return nil, err
	}

	var chunks []Chunk

	for _, document := range documents {
		docChunks, err := r.manager.LoadChunks(document.ID)
		if err != nil {
			return nil, err
		}

		chunks = append(chunks, docChunks...)
	}

	return chunks, nil
}
