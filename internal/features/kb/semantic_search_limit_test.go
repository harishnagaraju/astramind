package kb

import "testing"

// TestSemanticSearchCapsResults guards against a real bug found in
// manual testing: SemanticSearch previously returned every embedded
// chunk in the entire knowledge base with no limit, meaning a RAG
// prompt would eventually include the whole knowledge base on every
// question - diluting relevant content with irrelevant chunks, and
// eventually exceeding the model's context window outright as more
// documents get imported.
func TestSemanticSearchCapsResults(t *testing.T) {

	tempDir := t.TempDir()

	storage := NewJSONStorage(tempDir)
	manager := NewManager(storage)
	repository := NewRepository(manager)

	// Manually attach embeddings to more chunks than the limit
	// allows, so this test doesn't depend on a real embedder.
	total := DefaultSemanticSearchLimit + 3

	chunks := make([]Chunk, 0, total)

	for i := 0; i < total; i++ {
		chunks = append(chunks, Chunk{
			ID:         string(rune('a' + i)),
			DocumentID: "doc-1",
			Index:      i,
			Content:    "chunk content",
			// Slightly different vectors so ranking is deterministic
			// rather than all tied at the same score.
			Embedding: []float32{1, float32(i) * 0.01, 0},
		})
	}

	if err := manager.SaveDocument(&Document{
		ID:         "doc-1",
		Name:       "test.txt",
		ChunkCount: total,
	}); err != nil {
		t.Fatal(err)
	}

	if err := manager.SaveChunks(chunks); err != nil {
		t.Fatal(err)
	}

	results, err := repository.SemanticSearch([]float32{1, 0, 0})
	if err != nil {
		t.Fatal(err)
	}

	if len(results) != DefaultSemanticSearchLimit {
		t.Fatalf(
			"expected results capped at %d, got %d",
			DefaultSemanticSearchLimit,
			len(results),
		)
	}
}

func TestSemanticSearchDoesNotPadBelowLimit(t *testing.T) {

	tempDir := t.TempDir()

	storage := NewJSONStorage(tempDir)
	manager := NewManager(storage)
	repository := NewRepository(manager)

	chunks := []Chunk{
		{
			ID:         "only-one",
			DocumentID: "doc-1",
			Index:      0,
			Content:    "chunk content",
			Embedding:  []float32{1, 0, 0},
		},
	}

	if err := manager.SaveDocument(&Document{
		ID:         "doc-1",
		Name:       "test.txt",
		ChunkCount: 1,
	}); err != nil {
		t.Fatal(err)
	}

	if err := manager.SaveChunks(chunks); err != nil {
		t.Fatal(err)
	}

	results, err := repository.SemanticSearch([]float32{1, 0, 0})
	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 result when only 1 chunk exists, got %d", len(results))
	}
}
