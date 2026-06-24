package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	/*"github.com/harishnagaraju/astramind/internal/config"*/
	"github.com/harishnagaraju/astramind/internal/models"
)

func SessionExists(
    session string,
) bool {

    file := sessionFile(session)

    _, err := os.Stat(file)

    return err == nil
}

func DeleteSession(
    session string,
) error {

    file := sessionFile(session)

    return os.Remove(file)
}

func ListSessions() ([]string, error) {

    files, err := os.ReadDir("data/sessions")
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

        sessionName := strings.TrimSuffix(
            name,
            ".json",
        )

        sessions = append(
            sessions,
            sessionName,
        )
    }

    return sessions, nil
}

func CreateSession(session string) error {

    file := sessionFile(session)

    if _, err := os.Stat(file); err == nil {
        return nil
    }

    empty := []models.Message{}

    return SaveHistory(
        session,
        empty,
    )
}

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