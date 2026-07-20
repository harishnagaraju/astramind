package kb

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRepositoryChunks(t *testing.T) {
	tempDir := t.TempDir()

	storage := NewJSONStorage(tempDir)
	manager := NewManager(storage)
	repository := NewRepository(manager)

	source := filepath.Join(tempDir, "sample.txt")

	if err := os.WriteFile(source, []byte("Hello AstraMind"), 0644); err != nil {
		t.Fatal(err)
	}

	if _, err := manager.ImportDocument(source); err != nil {
		t.Fatal(err)
	}

	chunks, err := repository.Chunks()
	if err != nil {
		t.Fatal(err)
	}

	if len(chunks) == 0 {
		t.Fatal("expected at least one chunk")
	}
}
