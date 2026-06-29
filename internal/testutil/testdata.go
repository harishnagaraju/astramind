package testutil

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/harishnagaraju/astramind/internal/models"
)

func LoadConversation(
	t *testing.T,
	name string,
) []models.Message {

	t.Helper()

	root, err := os.Getwd()
	if err != nil {
		t.Fatalf("Unable to determine working directory: %v", err)
	}

	// Go back to project root
	root = filepath.Join(root, "..", "..")

	file := filepath.Join(
		root,
		"tests",
		"testdata",
		"conversations",
		name+".json",
	)

	data, err := os.ReadFile(file)
	if err != nil {
		t.Fatalf("Unable to read %s: %v", file, err)
	}

	var messages []models.Message

	err = json.Unmarshal(data, &messages)
	if err != nil {
		t.Fatalf("Invalid JSON in %s: %v", file, err)
	}

	return messages
}
