package history

import (
	"github.com/harishnagaraju/astramind/internal/infrastructure/models"
	"github.com/harishnagaraju/astramind/internal/infrastructure/storage"
)

// Service provides conversation history operations.
type Service struct {
	store Store
}

// NewService creates a history service backed by the real, on-disk
// data/sessions folder - the same location and behavior as before
// this store was made injectable. Every existing caller of
// NewService() is unaffected.
func NewService() *Service {
	return &Service{
		store: storage.NewFileHistoryStore("data"),
	}
}

// NewServiceWithStore creates a history service backed by the given
// store. Primarily for tests that want isolation from the real
// data/sessions folder - pass storage.NewFileHistoryStore(t.TempDir())
// for a store that's automatically cleaned up, even if the test
// fails partway through.
func NewServiceWithStore(store Store) *Service {
	return &Service{store: store}
}

// Save stores conversation history.
func (s *Service) Save(
	session string,
	messages []models.Message,
) error {

	return s.store.SaveHistory(session, messages)
}

// Load retrieves conversation history.
func (s *Service) Load(
	session string,
) ([]models.Message, error) {

	return s.store.LoadHistory(session)
}

// ListSessions returns the names of all saved sessions.
func (s *Service) ListSessions() ([]string, error) {

	return s.store.ListSessions()
}

// Delete removes a saved session.
func (s *Service) Delete(session string) error {

	return s.store.DeleteSession(session)
}
