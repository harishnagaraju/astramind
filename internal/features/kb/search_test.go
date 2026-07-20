package kb

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSearchSingleMatch(t *testing.T) {
	tempDir := t.TempDir()

	storage := NewJSONStorage(tempDir)
	manager := NewManager(storage)
	repository := NewRepository(manager)

	source := filepath.Join(tempDir, "sample.txt")

	text := "Go is a programming language."

	if err := os.WriteFile(source, []byte(text), 0644); err != nil {
		t.Fatal(err)
	}

	if _, err := manager.ImportDocument(source); err != nil {
		t.Fatal(err)
	}

	results, err := repository.Search("programming")
	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
}

func TestSearchNoMatch(t *testing.T) {
	tempDir := t.TempDir()

	storage := NewJSONStorage(tempDir)
	manager := NewManager(storage)
	repository := NewRepository(manager)

	source := filepath.Join(tempDir, "sample.txt")

	if err := os.WriteFile(source, []byte("Hello World"), 0644); err != nil {
		t.Fatal(err)
	}

	if _, err := manager.ImportDocument(source); err != nil {
		t.Fatal(err)
	}

	results, err := repository.Search("golang")
	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

func TestSearchCaseInsensitive(t *testing.T) {
	tempDir := t.TempDir()

	storage := NewJSONStorage(tempDir)
	manager := NewManager(storage)
	repository := NewRepository(manager)

	source := filepath.Join(tempDir, "sample.txt")

	if err := os.WriteFile(source, []byte("OpenAI develops AI systems."), 0644); err != nil {
		t.Fatal(err)
	}

	if _, err := manager.ImportDocument(source); err != nil {
		t.Fatal(err)
	}

	results, err := repository.Search("openai")
	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
}

func TestSearchEmptyQuery(t *testing.T) {
	tempDir := t.TempDir()

	storage := NewJSONStorage(tempDir)
	manager := NewManager(storage)
	repository := NewRepository(manager)

	results, err := repository.Search("")
	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}
