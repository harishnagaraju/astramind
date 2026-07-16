package session

import (
	"github.com/harishnagaraju/astramind/internal/infrastructure/models"
	"github.com/harishnagaraju/astramind/internal/infrastructure/storage"
)

// Service provides session management operations.
type Service struct{}

// NewService creates a new session service.
func NewService() *Service {
	return &Service{}
}

// Save stores a conversation as a session.
func (s *Service) Save(
	name string,
	messages []models.Message,
) error {

	return storage.SaveHistory(name, messages)
}

// Load retrieves a saved session.
func (s *Service) Load(
	name string,
) ([]models.Message, error) {

	return storage.LoadHistory(name)
}

// Delete removes a saved session.
func (s *Service) Delete(
	name string,
) error {

	return storage.DeleteSession(name)
}

// List returns all available sessions.
func (s *Service) List() ([]string, error) {

	return storage.ListSessions()
}

// Create creates a new session.
func (s *Service) Create(
	name string,
) error {

	return storage.CreateSession(name)
}

// Exists checks whether a session exists.
func (s *Service) Exists(
	name string,
) bool {

	return storage.SessionExists(name)
}
