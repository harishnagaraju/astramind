package kb

import (
	"strings"
	"testing"
)

// TestFilterRelevantChunks_Issue1RealFailureCase reproduces the real
// bug found via manual testing on real Ollama: a knowledge base
// containing both the real Sanskrit1.txt content and an unrelated
// service contract answered a Sanskrit-specific enumeration question
// by bundling both documents together, because SemanticSearch
// returns every chunk whenever the KB's total chunk count is at or
// below its retrieval limit (5) - true for any small KB regardless
// of embedding quality.
func TestFilterRelevantChunks_Issue1RealFailureCase(t *testing.T) {
	doc := &Document{ID: "sanskrit-doc", Content: realSanskritNoticeboard}
	sanskritChunks := ChunkDocument(doc, DefaultChunkSize, DefaultOverlap)

	contractDoc := &Document{
		ID: "contract-doc",
		Content: "SERVICE AGREEMENT\n\n" +
			"This Service Agreement is entered into as of January 15, 2026, between Acme Consulting LLC and Beta Client Corp.\n\n" +
			"Client shall pay Consultant a fixed fee of $45,000, payable in three installments of $15,000 each.",
	}
	contractChunks := ChunkDocument(contractDoc, DefaultChunkSize, DefaultOverlap)

	var results []SemanticSearchResult
	for _, c := range sanskritChunks {
		results = append(results, SemanticSearchResult{
			DocumentID: c.DocumentID, ChunkID: c.ID, ChunkIndex: c.Index, Content: c.Content,
		})
	}
	for _, c := range contractChunks {
		results = append(results, SemanticSearchResult{
			DocumentID: c.DocumentID, ChunkID: c.ID, ChunkIndex: c.Index, Content: c.Content,
		})
	}

	filtered := FilterRelevantChunks("what is the timings of sanskrit classes?", results)

	if len(filtered) != len(sanskritChunks) {
		t.Fatalf("expected exactly the %d Sanskrit chunks, got %d chunks", len(sanskritChunks), len(filtered))
	}

	for _, r := range filtered {
		if r.DocumentID == "contract-doc" {
			t.Fatal("expected the unrelated contract to be filtered out, but it was included")
		}
	}
}

// TestFilterRelevantChunks_DoesNotDropRelatedEntriesWithinAChunk is
// the critical regression guard: filtering must operate at whole-
// chunk granularity only. A chunk containing "Tuesday Chanting"
// alongside genuinely Sanskrit-related entries must be kept whole -
// this is exactly the failure mode that item-level filtering caused
// and was rejected for earlier in this project's history.
func TestFilterRelevantChunks_DoesNotDropRelatedEntriesWithinAChunk(t *testing.T) {
	doc := &Document{ID: "sanskrit-doc", Content: realSanskritNoticeboard}
	chunks := ChunkDocument(doc, DefaultChunkSize, DefaultOverlap)

	var results []SemanticSearchResult
	for _, c := range chunks {
		results = append(results, SemanticSearchResult{
			DocumentID: c.DocumentID, ChunkID: c.ID, ChunkIndex: c.Index, Content: c.Content,
		})
	}

	filtered := FilterRelevantChunks("what is the timings of sanskrit classes?", results)

	// Every chunk that survives filtering must retain its full
	// original content - "Tuesday Chanting" (which doesn't contain
	// the word "Sanskrit") must still be present if it shares a
	// chunk with entries that do.
	foundChanting := false
	for _, r := range filtered {
		if strings.Contains(r.Content, "Chanting") {
			foundChanting = true
		}
	}
	if !foundChanting {
		t.Fatal("expected 'Chanting' entries to survive chunk-level filtering (they must not be individually dropped)")
	}
}

// TestFilterRelevantChunks_AllChunksFilteredFallsBackToUnfiltered
// confirms the safety fallback: if filtering would remove every
// chunk, return the original list rather than nothing at all.
func TestFilterRelevantChunks_AllChunksFilteredFallsBackToUnfiltered(t *testing.T) {
	results := []SemanticSearchResult{
		{DocumentID: "d1", ChunkID: "c1", Content: "Completely unrelated content about gardening."},
	}

	filtered := FilterRelevantChunks("what is the zoom meeting id", results)

	if len(filtered) != 1 {
		t.Fatalf("expected fallback to the original unfiltered list, got %d results", len(filtered))
	}
}

func TestFilterRelevantChunks_EmptyQuestion(t *testing.T) {
	results := []SemanticSearchResult{
		{DocumentID: "d1", ChunkID: "c1", Content: "Some content."},
	}

	filtered := FilterRelevantChunks("", results)

	if len(filtered) != 1 {
		t.Fatal("expected an empty/stopword-only question to pass everything through unfiltered")
	}
}