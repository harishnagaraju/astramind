package kb

import "testing"

// TestHasAnyContentWordOverlap_RealFailureCase reproduces the exact
// bug found in real manual testing: "what is the capital of delhi"
// (a question with zero relevance to the knowledge base) was
// confidently answered using an entire unrelated chunk of Sanskrit
// class schedule content, because SemanticSearch always returns its
// top-ranked chunks with no minimum relevance check at all.
func TestHasAnyContentWordOverlap_RealFailureCase(t *testing.T) {
	doc := &Document{ID: "sanskrit-doc", Content: realSanskritNoticeboard}
	chunks := ChunkDocument(doc, DefaultChunkSize, DefaultOverlap)

	var results []SemanticSearchResult
	for _, c := range chunks {
		results = append(results, SemanticSearchResult{
			DocumentID: c.DocumentID, ChunkID: c.ID, ChunkIndex: c.Index, Content: c.Content,
		})
	}

	// The real out-of-scope question from manual testing - shares no
	// words with anything in the Sanskrit noticeboard.
	if HasAnyContentWordOverlap("what is the capital of delhi", results) {
		t.Fatal("expected no overlap for a completely unrelated question, got overlap")
	}

	// Confirm a genuinely relevant question is NOT blocked by the
	// same gate - this must not regress real, working questions.
	if !HasAnyContentWordOverlap("what are the sanskrit class timings", results) {
		t.Fatal("expected overlap for a genuinely relevant question, got none")
	}

	if !HasAnyContentWordOverlap("what is the zoom meeting id", results) {
		t.Fatal("expected overlap for a genuinely relevant single-fact question, got none")
	}
}

func TestHasAnyContentWordOverlap_EmptyQuestion(t *testing.T) {
	results := []SemanticSearchResult{{DocumentID: "d", ChunkID: "c", Content: "Some content."}}

	// A question with no content words at all (only stopwords)
	// shouldn't be blocked - there's nothing meaningful to check
	// overlap against, so err on the side of not gating.
	if !HasAnyContentWordOverlap("what is the", results) {
		t.Fatal("expected a stopword-only question to pass through ungated")
	}
}

func TestHasAnyContentWordOverlap_NoResults(t *testing.T) {
	if HasAnyContentWordOverlap("what are the sanskrit class timings", nil) {
		t.Fatal("expected no overlap when there are no results at all")
	}
}

// TestHasAnyContentWordOverlap_ShortWordFalsePositive is a real
// false-positive found in live testing against the actual binary,
// not a hypothetical: "id" (2 characters) matched as a literal
// substring inside "Confidentiality", letting a completely unrelated
// document (a service contract with no Zoom/meeting content at all)
// pass the relevance gate for "what is the zoom meeting id".
func TestHasAnyContentWordOverlap_ShortWordFalsePositive(t *testing.T) {
	results := []SemanticSearchResult{
		{
			DocumentID: "contract",
			ChunkID:    "c1",
			Content:    "Both parties agree to maintain strict confidentiality regarding all proprietary information.",
		},
	}

	if HasAnyContentWordOverlap("what is the zoom meeting id", results) {
		t.Fatal("expected no overlap - 'id' matching inside 'confidentiality' is a false positive, not real relevance")
	}

	// Confirm the real content words ("zoom", "meeting") still work
	// correctly when they're genuinely present - this fix should
	// only exclude very short words, not break real matching.
	relevantResults := []SemanticSearchResult{
		{DocumentID: "d", ChunkID: "c", Content: "Meeting ID 795 777 3585, hosted via Zoom."},
	}
	if !HasAnyContentWordOverlap("what is the zoom meeting id", relevantResults) {
		t.Fatal("expected overlap when zoom/meeting are genuinely present")
	}
}
