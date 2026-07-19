package storage

import (
	"github.com/harishnagaraju/astramind/internal/infrastructure/models"
)

// defaultHistoryStore is the store backing the package-level
// functions below, hardcoded to "data" for backward compatibility
// with every existing caller. Code that wants an isolated store
// (tests, primarily) should use NewFileHistoryStore directly instead
// of these package-level functions.
var defaultHistoryStore = NewFileHistoryStore("data")

func SessionExists(session string) bool {
	return defaultHistoryStore.SessionExists(session)
}

func DeleteSession(session string) error {
	return defaultHistoryStore.DeleteSession(session)
}

func ListSessions() ([]string, error) {
	return defaultHistoryStore.ListSessions()
}

func CreateSession(session string) error {
	return defaultHistoryStore.CreateSession(session)
}

func LoadHistory(session string) ([]models.Message, error) {
	return defaultHistoryStore.LoadHistory(session)
}

func SaveHistory(
	session string,
	messages []models.Message,
) error {
	return defaultHistoryStore.SaveHistory(session, messages)
}
