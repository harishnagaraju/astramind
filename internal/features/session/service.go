package session

import (
	"github.com/harishnagaraju/astramind/internal/features/history"
	"github.com/harishnagaraju/astramind/internal/infrastructure/models"
	"github.com/harishnagaraju/astramind/internal/infrastructure/storage"
)

// Service provides session management operations. It delegates
// content persistence (Save/Load/Delete/List) to the History feature,
// since a "session" is a named history record - and adds session
// lifecycle operations (Create/Exists) that don't belong in History.
type Service struct {
	history *history.Service
	store   Store
}

// NewService creates a new session service backed by the given
// History service. Create/Exists default to the real, on-disk
// data/sessions folder - the same location and behavior as before
// this store was made injectable. Every existing caller of
// NewService() is unaffected.
func NewService(historySvc *history.Service) *Service {
	return &Service{
		history: historySvc,
		store:   storage.NewFileHistoryStore("data"),
	}
}

// NewServiceWithStore creates a session service backed by the given
// History service and Store. Primarily for tests that want isolation
// from the real data/sessions folder.
func NewServiceWithStore(historySvc *history.Service, store Store) *Service {
	return &Service{
		history: historySvc,
		store:   store,
	}
}

// Save stores a conversation as a session.
func (s *Service) Save(
	name string,
	messages []models.Message,
) error {

	return s.history.Save(name, messages)
}

// Load retrieves a saved session.
func (s *Service) Load(
	name string,
) ([]models.Message, error) {

	return s.history.Load(name)
}

// Delete removes a saved session.
func (s *Service) Delete(
	name string,
) error {

	return s.history.Delete(name)
}

// List returns all available sessions.
func (s *Service) List() ([]string, error) {

	return s.history.ListSessions()
}

// Create creates a new session.
func (s *Service) Create(
	name string,
) error {

	return s.store.CreateSession(name)
}

// Exists checks whether a session exists.
func (s *Service) Exists(
	name string,
) bool {

	return s.store.SessionExists(name)
}
