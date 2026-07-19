package history

import "github.com/harishnagaraju/astramind/internal/infrastructure/models"

// Store is what Service needs to persist session history. It is
// satisfied structurally by *storage.FileHistoryStore - defined here
// so history does not need to depend on the concrete storage type,
// and so tests can substitute a store rooted at a temp directory
// instead of the real data/sessions folder.
type Store interface {
	SaveHistory(session string, messages []models.Message) error
	LoadHistory(session string) ([]models.Message, error)
	DeleteSession(session string) error
	ListSessions() ([]string, error)
}
