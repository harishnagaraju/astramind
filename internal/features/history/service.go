package history

import (
	"github.com/harishnagaraju/astramind/internal/infrastructure/models"
	"github.com/harishnagaraju/astramind/internal/infrastructure/storage"
)

// Service provides conversation history operations.
type Service struct{}

// NewService creates a history service.
func NewService() *Service {
	return &Service{}
}

// Save stores conversation history.
func (s *Service) Save(
	session string,
	messages []models.Message,
) error {

	return storage.SaveHistory(session, messages)
}

// Load retrieves conversation history.
func (s *Service) Load(
	session string,
) ([]models.Message, error) {

	return storage.LoadHistory(session)
}

// ListSessions returns the names of all saved sessions.
func (s *Service) ListSessions() ([]string, error) {

	return storage.ListSessions()
}

// Delete removes a saved session.
func (s *Service) Delete(session string) error {
	return storage.DeleteSession(session)
}
