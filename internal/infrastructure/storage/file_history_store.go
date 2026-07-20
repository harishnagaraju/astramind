package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/harishnagaraju/astramind/internal/infrastructure/models"
)

// FileHistoryStore persists session history to JSON files under a
// configurable root directory. Unlike the package-level functions in
// history.go (which are hardcoded to "data" for backward
// compatibility), this type can be pointed at a temp directory in
// tests, so tests never touch the real data/sessions folder.
type FileHistoryStore struct {
	directory string
}

// NewFileHistoryStore creates a store rooted at the given directory.
func NewFileHistoryStore(directory string) *FileHistoryStore {
	return &FileHistoryStore{directory: directory}
}

func (s *FileHistoryStore) sessionFile(session string) string {
	return filepath.Join(
		s.directory,
		"sessions",
		session+".json",
	)
}

func (s *FileHistoryStore) SessionExists(session string) bool {
	_, err := os.Stat(s.sessionFile(session))
	return err == nil
}

func (s *FileHistoryStore) DeleteSession(session string) error {
	return os.Remove(s.sessionFile(session))
}

func (s *FileHistoryStore) ListSessions() ([]string, error) {

	files, err := os.ReadDir(
		filepath.Join(s.directory, "sessions"),
	)
	if err != nil {
		return nil, err
	}

	sessions := []string{}

	for _, file := range files {

		if file.IsDir() {
			continue
		}

		name := file.Name()

		if filepath.Ext(name) != ".json" {
			continue
		}

		sessions = append(
			sessions,
			strings.TrimSuffix(name, ".json"),
		)
	}

	return sessions, nil
}

func (s *FileHistoryStore) CreateSession(session string) error {

	file := s.sessionFile(session)

	if _, err := os.Stat(file); err == nil {
		return nil
	}

	return s.SaveHistory(session, []models.Message{})
}

func (s *FileHistoryStore) LoadHistory(
	session string,
) ([]models.Message, error) {

	file := s.sessionFile(session)

	if _, err := os.Stat(file); os.IsNotExist(err) {
		return []models.Message{}, nil
	}

	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var messages []models.Message

	if err := json.Unmarshal(data, &messages); err != nil {
		return nil, err
	}

	return messages, nil
}

func (s *FileHistoryStore) SaveHistory(
	session string,
	messages []models.Message,
) error {

	file := s.sessionFile(session)

	if err := os.MkdirAll(filepath.Dir(file), 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(messages, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(file, data, 0644)
}
