package search

import (
	"github.com/harishnagaraju/astramind/internal/features/history"
	"github.com/harishnagaraju/astramind/internal/infrastructure/models"
)

// Service provides search operations. It depends on the History feature,
// not on infrastructure/storage directly.
type Service struct {
	history *history.Service
}

// NewService creates a search service backed by the given History service.
func NewService(historySvc *history.Service) *Service {
	return &Service{history: historySvc}
}

// SearchCurrent searches an already-loaded conversation in memory.
func (s *Service) SearchCurrent(
	messages []models.Message,
	query string,
) []models.SearchResult {

	return SearchMessages(messages, query)
}

// SearchAll searches every saved session via the History feature.
func (s *Service) SearchAll(
	query string,
) ([]models.SessionSearchResult, error) {

	return SearchAllSessions(s.history, query)
}
