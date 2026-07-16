package kb

import (
	"strings"
	"testing"
)

func TestChunkSmallDocument(t *testing.T) {
	doc := &Document{
		ID:      "doc1",
		Content: "Hello World",
	}

	chunks := ChunkDocument(doc, 100, 10)

	if len(chunks) != 1 {
		t.Fatalf("expected 1 chunk, got %d", len(chunks))
	}

	if chunks[0].DocumentID != doc.ID {
		t.Fatalf("expected document ID %q, got %q", doc.ID, chunks[0].DocumentID)
	}

	if chunks[0].Index != 0 {
		t.Fatalf("expected chunk index 0, got %d", chunks[0].Index)
	}

	if chunks[0].Content != doc.Content {
		t.Fatal("chunk content does not match original document")
	}
}

func TestChunkLargeDocument(t *testing.T) {
	content := strings.Repeat("A", 5000)

	doc := &Document{
		ID:      "doc1",
		Content: content,
	}

	chunks := ChunkDocument(doc, 1000, 100)

	if len(chunks) <= 1 {
		t.Fatalf("expected multiple chunks, got %d", len(chunks))
	}

	for i, chunk := range chunks {
		if chunk.Index != i {
			t.Fatalf("expected chunk index %d, got %d", i, chunk.Index)
		}

		if chunk.DocumentID != doc.ID {
			t.Fatalf("expected document ID %q, got %q", doc.ID, chunk.DocumentID)
		}

		if len(chunk.Content) == 0 {
			t.Fatalf("chunk %d is empty", i)
		}
	}
}

func TestChunkOverlap(t *testing.T) {
	// 260 characters: ABCDEFGHIJKLMNOPQRSTUVWXYZ repeated 10 times.
	content := strings.Repeat("ABCDEFGHIJKLMNOPQRSTUVWXYZ", 10)

	doc := &Document{
		ID:      "doc1",
		Content: content,
	}

	chunkSize := 50
	overlap := 10

	chunks := ChunkDocument(doc, chunkSize, overlap)

	if len(chunks) < 2 {
		t.Fatalf("expected multiple chunks, got %d", len(chunks))
	}

	// Verify overlap mathematically.
	expectedOverlap := chunks[0].Content[chunkSize-overlap:]

	actualOverlap := chunks[1].Content[:overlap]

	if expectedOverlap != actualOverlap {
		t.Fatalf("expected overlap %q, got %q", expectedOverlap, actualOverlap)
	}

	for i, chunk := range chunks {
		if chunk.Index != i {
			t.Fatalf("expected chunk index %d, got %d", i, chunk.Index)
		}
	}
}
