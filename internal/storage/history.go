package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	/*"github.com/harishnagaraju/astramind/internal/config"*/
	"github.com/harishnagaraju/astramind/internal/models"
)

func sessionFile(session string) string {

	return filepath.Join(
		"data",
		"sessions",
		session+".json",
	)
}

func LoadHistory(
	session string,
) ([]models.Message, error) {

	file := sessionFile(session)

	if _, err := os.Stat(file); os.IsNotExist(err) {
		return []models.Message{}, nil
	}

	data, err := os.ReadFile(file)

	if err != nil {
		return nil, err
	}

	var messages []models.Message

	err = json.Unmarshal(
		data,
		&messages,
	)

	if err != nil {
		return nil, err
	}

	return messages, nil
}

func SaveHistory(
	session string,
	messages []models.Message,
) error {

	file := sessionFile(session)

	err := os.MkdirAll(
		filepath.Dir(file),
		0755,
	)

	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(
		messages,
		"",
		"  ",
	)

	if err != nil {
		return err
	}

	return os.WriteFile(
		file,
		data,
		0644,
	)
}