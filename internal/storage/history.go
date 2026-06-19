package storage

import (
	"encoding/json"
	"os"
	"github.com/harishnagaraju/astramind/internal/models"
)

func LoadHistory() ([]models.Message, error) {

	data, err := os.ReadFile("data/chat_history.json")

	if err != nil {

		if os.IsNotExist(err) {
			return []models.Message{}, nil
		}

		return nil, err
	}

	var messages []models.Message

	err = json.Unmarshal(data, &messages)

	if err != nil {
		return nil, err
	}

	return messages, nil

}

func SaveHistory(messages []models.Message) error {

	data, err := json.MarshalIndent(
		messages,
		"",
		"  ",
	)

	if err != nil {
		return err
	}

	return os.WriteFile(
		"data/chat_history.json",
		data,
		0644,
	)

}
