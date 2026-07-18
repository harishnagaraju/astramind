package kb

// Service provides Knowledge Base operations.
type Service struct {
	manager *Manager
}

// NewService creates a KB service.
func NewService(manager *Manager) *Service {
	return &Service{
		manager: manager,
	}
}

// Import imports a document into the knowledge base.
func (s *Service) Import(path string) (*Document, error) {
	return s.manager.ImportDocument(path)
}

// List returns all imported documents.
func (s *Service) List() ([]Document, error) {
	return s.manager.ListKnowledge()
}

// Search searches the knowledge base.
func (s *Service) Search(query string) ([]SearchResult, error) {
	return s.manager.Search(query)
}

// SemanticSearch performs an embedding-based search of the knowledge base.
func (s *Service) SemanticSearch(query string) ([]SemanticSearchResult, error) {
	return s.manager.SemanticSearch(query)
}

// Remove removes a document.
func (s *Service) Remove(id string) error {
	return s.manager.RemoveKnowledge(id)
}

// Clear removes all documents.
func (s *Service) Clear() error {
	return s.manager.ClearKnowledge()
}

// Stats returns KB statistics.
func (s *Service) Stats() (*Stats, error) {
	return s.manager.Stats()
}
