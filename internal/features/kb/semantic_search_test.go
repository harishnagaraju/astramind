package kb

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

// fakeEmbedder is a test double for Embedder. It returns a fixed
// vector for every call, or an error if failOnEmbed is set.
type fakeEmbedder struct {
	failOnEmbed bool
	callCount   int
}

func (f *fakeEmbedder) Embed(text string) ([]float32, error) {

	f.callCount++

	if f.failOnEmbed {
		return nil, errors.New("embedding failed")
	}

	return []float32{0.1, 0.2, 0.3}, nil
}

func TestImportDocument_WithEmbedder(t *testing.T) {
	tempDir := t.TempDir()

	source := filepath.Join(tempDir, "sample.txt")

	err := os.WriteFile(source, []byte("Hello AstraMind"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	storage := NewJSONStorage(filepath.Join(tempDir, "kb"))
	manager := NewManager(storage)

	embedder := &fakeEmbedder{}
	manager.SetEmbedder(embedder)

	doc, err := manager.ImportDocument(source)
	if err != nil {
		t.Fatal(err)
	}

	if embedder.callCount == 0 {
		t.Fatal("expected embedder to be called during import")
	}

	chunks, err := manager.LoadChunks(doc.ID)
	if err != nil {
		t.Fatal(err)
	}

	if len(chunks) == 0 {
		t.Fatal("expected at least one chunk")
	}

	for _, chunk := range chunks {
		if len(chunk.Embedding) == 0 {
			t.Fatal("expected chunk to have an embedding")
		}
	}
}

func TestImportDocument_EmbedderFailureDoesNotFailImport(t *testing.T) {
	tempDir := t.TempDir()

	source := filepath.Join(tempDir, "sample.txt")

	err := os.WriteFile(source, []byte("Hello AstraMind"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	storage := NewJSONStorage(filepath.Join(tempDir, "kb"))
	manager := NewManager(storage)

	embedder := &fakeEmbedder{failOnEmbed: true}
	manager.SetEmbedder(embedder)

	doc, err := manager.ImportDocument(source)
	if err != nil {
		t.Fatalf("import should succeed even if embedding fails, got: %v", err)
	}

	chunks, err := manager.LoadChunks(doc.ID)
	if err != nil {
		t.Fatal(err)
	}

	for _, chunk := range chunks {
		if len(chunk.Embedding) != 0 {
			t.Fatal("expected no embedding when embedder fails")
		}
	}
}

func TestImportDocument_NoEmbedderConfigured(t *testing.T) {
	tempDir := t.TempDir()

	source := filepath.Join(tempDir, "sample.txt")

	err := os.WriteFile(source, []byte("Hello AstraMind"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	storage := NewJSONStorage(filepath.Join(tempDir, "kb"))
	manager := NewManager(storage)

	// No SetEmbedder call - matches existing pre-embedding behavior.
	doc, err := manager.ImportDocument(source)
	if err != nil {
		t.Fatal(err)
	}

	chunks, err := manager.LoadChunks(doc.ID)
	if err != nil {
		t.Fatal(err)
	}

	for _, chunk := range chunks {
		if len(chunk.Embedding) != 0 {
			t.Fatal("expected no embedding when no embedder is configured")
		}
	}
}
