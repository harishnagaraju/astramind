package kb

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRemoveKnowledge(t *testing.T) {
	tempDir := t.TempDir()

	storage := NewJSONStorage(tempDir)
	manager := NewManager(storage)

	source := filepath.Join(tempDir, "sample.txt")

	if err := os.WriteFile(source, []byte("Hello AstraMind"), 0644); err != nil {
		t.Fatal(err)
	}

	doc, err := manager.ImportDocument(source)
	if err != nil {
		t.Fatal(err)
	}

	if err := manager.RemoveKnowledge(doc.ID); err != nil {
		t.Fatal(err)
	}

	_, err = manager.GetKnowledge(doc.ID)
	if err == nil {
		t.Fatal("expected document to be removed")
	}
}

func TestClearKnowledge(t *testing.T) {
	tempDir := t.TempDir()

	storage := NewJSONStorage(tempDir)
	manager := NewManager(storage)

	files := []string{
		"a.txt",
		"b.txt",
		"c.txt",
	}

	for _, file := range files {
		path := filepath.Join(tempDir, file)

		if err := os.WriteFile(path, []byte("Knowledge"), 0644); err != nil {
			t.Fatal(err)
		}

		if _, err := manager.ImportDocument(path); err != nil {
			t.Fatal(err)
		}
	}

	if err := manager.ClearKnowledge(); err != nil {
		t.Fatal(err)
	}

	documents, err := manager.ListKnowledge()
	if err != nil {
		t.Fatal(err)
	}

	if len(documents) != 0 {
		t.Fatalf("expected empty knowledge base, got %d documents", len(documents))
	}
}
