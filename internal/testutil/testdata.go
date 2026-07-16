package testutil

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/harishnagaraju/astramind/internal/infrastructure/models"
)

func LoadConversation(
	t *testing.T,
	name string,
) []models.Message {

	t.Helper()

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("failed to determine testutil source location")
	}

	// testdata.go -> internal/testutil -> project root
	root := filepath.Join(
		filepath.Dir(filename),
		"..",
		"..",
	)

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
