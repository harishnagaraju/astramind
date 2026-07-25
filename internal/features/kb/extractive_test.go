package kb

import (
	"strings"
	"testing"
)

// bagOfWordsEmbedder is a small, deterministic, dependency-free
// embedder for testing ranking logic without a live model. See
// manager.go's ExtractiveAnswer doc comment for why a windowed
// version of this function was tried and abandoned: real embedding
// similarity data, captured live against Ollama, showed no usable
// topic-boundary signal at the sentence level for this content - an
// unrelated sentence scored HIGHER against the anchor than genuinely
// related sentences did, apparently due to incidental word overlap
// rather than real relevance. That finding came from real model
// output, not from this fake, which only tests literal word overlap
// and is used here purely to verify chunk-selection ranking.
type bagOfWordsEmbedder struct {
	vocabulary []string
}

func (e *bagOfWordsEmbedder) Embed(text string) ([]float32, error) {
	lower := strings.ToLower(text)
	vec := make([]float32, len(e.vocabulary))
	for i, word := range e.vocabulary {
		if strings.Contains(lower, word) {
			vec[i] = 1
		}
	}
	return vec, nil
}

// TestExtractiveAnswer_SelectsCorrectChunk confirms ExtractiveAnswer
// picks the chunk containing the best-matching sentence, and returns
// that chunk's full content - not a window, not a single sentence.
func TestExtractiveAnswer_SelectsCorrectChunk(t *testing.T) {
	doc := &Document{ID: "sanskrit-doc", Content: realSanskritNoticeboard}
	chunks := ChunkDocument(doc, DefaultChunkSize, DefaultOverlap)

	var results []SemanticSearchResult
	for _, c := range chunks {
		results = append(results, SemanticSearchResult{
			DocumentID: c.DocumentID, ChunkID: c.ID, ChunkIndex: c.Index, Content: c.Content,
		})
	}

	embedder := &bagOfWordsEmbedder{
		vocabulary: []string{
			"thursday", "senior", "monday", "youth", "tuesday",
			"chanting", "wednesday", "friday", "saturday", "evenings",
		},
	}
	manager := &Manager{}
	manager.SetEmbedder(embedder)

	item, err := manager.ExtractiveAnswer(
		"what time is thursday senior sanskrit",
		results,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(item.Text, "Thursday Senior Sanskrit from 06:15 to 07:15 am") {
		t.Fatalf("expected the matched fact in the returned text, got %q", item.Text)
	}

	// Precise line-matching (see precise_match.go) now correctly
	// narrows this to just the matched entry, not the whole chunk -
	// a real precision improvement verified against the exact real
	// case (this same question) that motivated it. The old
	// assertion here (checking for a full-chunk match) reflected the
	// previous, coarser whole-chunk-only behavior and no longer
	// describes correct behavior.
	if item.Text != "Thursday Senior Sanskrit from 06:15 to 07:15 am." {
		t.Fatalf("expected precise single-line match, got %q", item.Text)
	}
}

func TestExtractiveAnswer_NoEmbedderConfigured(t *testing.T) {
	manager := &Manager{}

	results := []SemanticSearchResult{
		{DocumentID: "doc1", ChunkID: "chunk1", Content: "Some content."},
	}

	_, err := manager.ExtractiveAnswer("a question", results)
	if err == nil {
		t.Fatal("expected error when no embedder is configured")
	}
}

func TestExtractiveAnswer_NoResults(t *testing.T) {
	embedder := &bagOfWordsEmbedder{vocabulary: []string{"test"}}
	manager := &Manager{}
	manager.SetEmbedder(embedder)

	_, err := manager.ExtractiveAnswer("a question", nil)
	if err != ErrNoExtractiveMatch {
		t.Fatalf("expected ErrNoExtractiveMatch, got %v", err)
	}
}
