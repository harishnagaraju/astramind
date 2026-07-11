package kb

// Storage defines the interface for persisting knowledge base data.
type Storage interface {
	SaveDocument(doc *Document) error
	LoadDocument(id string) (*Document, error)
	DeleteDocument(id string) error
	ListDocuments() ([]Document, error)

	SaveChunks(chunks []Chunk) error
	LoadChunks(documentID string) ([]Chunk, error)
	DeleteChunks(documentID string) error
}
