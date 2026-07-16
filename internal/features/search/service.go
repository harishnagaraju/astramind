package search

import (
	"github.com/harishnagaraju/astramind/internal/infrastructure/models"
	"github.com/harishnagaraju/astramind/internal/infrastructure/storage"
)

// Service provides search operations.
type Service struct{}

// NewService creates a search service.
func NewService() *Service {
	return &Service{}
}

// SearchCurrent searches the active conversation.
func (s *Service) SearchCurrent(
	messages []models.Message,
	query string,
) []models.SearchResult {

	return storage.SearchMessages(
		messages,
		query,
	)
}

// SearchAll searches every saved session.
func (s *Service) SearchAll(
	query string,
) ([]models.SessionSearchResult, error) {

	return storage.SearchAllSessions(query)
}
