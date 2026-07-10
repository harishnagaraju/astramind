package kb

// Storage defines the interface for persisting knowledge base documents.
type Storage interface {
	SaveDocument(doc *Document) error
	LoadDocument(id string) (*Document, error)
	DeleteDocument(id string) error
	ListDocuments() ([]Document, error)
}
