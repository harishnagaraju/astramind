package kb

import (
	"os"
	"path/filepath"
	"testing"
)

func TestImportDocument(t *testing.T) {
	tempDir := t.TempDir()

	source := filepath.Join(tempDir, "sample.txt")

	err := os.WriteFile(source, []byte("Hello AstraMind"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	storage := NewJSONStorage(filepath.Join(tempDir, "kb"))
	manager := NewManager(storage)

	doc, err := manager.ImportDocument(source)
	if err != nil {
		t.Fatal(err)
	}

	if doc == nil {
		t.Fatal("expected imported document")
	}

	if doc.Name != "sample.txt" {
		t.Fatalf("expected sample.txt, got %s", doc.Name)
	}

	if doc.Content != "Hello AstraMind" {
		t.Fatalf("unexpected content")
	}
}

func TestImportUnsupportedFile(t *testing.T) {
	tempDir := t.TempDir()

	source := filepath.Join(tempDir, "sample.pdf")

	err := os.WriteFile(source, []byte("dummy"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	storage := NewJSONStorage(filepath.Join(tempDir, "kb"))
	manager := NewManager(storage)

	_, err = manager.ImportDocument(source)
	if err == nil {
		t.Fatal("expected error")
	}
}
