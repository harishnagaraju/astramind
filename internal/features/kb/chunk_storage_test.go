package kb

import (
	"testing"
)

func TestSaveAndLoadChunks(t *testing.T) {
	dir := t.TempDir()

	storage := NewJSONStorage(dir)

	chunks := []Chunk{
		{
			ID:         "1",
			DocumentID: "doc1",
			Index:      0,
			Content:    "Hello",
		},
		{
			ID:         "2",
			DocumentID: "doc1",
			Index:      1,
			Content:    "World",
		},
	}

	if err := storage.SaveChunks(chunks); err != nil {
		t.Fatal(err)
	}

	loaded, err := storage.LoadChunks("doc1")
	if err != nil {
		t.Fatal(err)
	}

	if len(loaded) != 2 {
		t.Fatalf("expected 2 chunks, got %d", len(loaded))
	}
}
