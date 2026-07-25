package kb

import "strings"

// stopWords are common English function words stripped from a
// question before precise-line matching, so scoring focuses on the
// actual content words (names, topics, specific terms) rather than
// every word in the question contributing equally.
var stopWords = map[string]bool{
	"a": true, "an": true, "the": true, "is": true, "are": true,
	"was": true, "were": true, "what": true, "who": true, "when": true,
	"where": true, "how": true, "does": true, "do": true, "did": true,
	"for": true, "of": true, "to": true, "in": true, "on": true,
	"at": true, "there": true, "it": true, "and": true, "or": true,
	"with": true, "about": true,
}

// contentWords extracts the meaningful words from a question,
// stripping stopwords and punctuation. Order and duplicates are not
// significant here - this feeds a simple overlap count, not a
// phrase match.
func contentWords(question string) []string {
	var out []string
	for _, w := range strings.Fields(strings.ToLower(question)) {
		w = strings.Trim(w, ".,?!:;\"'")
		if w != "" && !stopWords[w] {
			out = append(out, w)
		}
	}
	return out
}

// findPreciseLine looks within a single chunk's content for one
// paragraph (per splitParagraphs) that contains strictly more of the
// question's content words than every other paragraph in the same
// chunk. Returns that paragraph and true only when there is a single,
// unambiguous best match - not a tie, and not a zero score.
//
// This exists as a precision refinement on top of ExtractiveAnswer's
// existing whole-chunk return, not a replacement for it. A prior
// windowing approach based on embedding cosine similarity was tried
// and abandoned after real measured data showed an unrelated sentence
// ("Not meeting on 16 February") scoring HIGHER similarity to a query
// ("Meeting ID 795 777 3585") than the genuinely relevant sentence
// did, apparently from incidental semantic proximity rather than real
// relevance - no threshold could separate the two reliably.
//
// Simple content-word overlap counting does not have that failure
// mode for this case: "Not meeting on 16 February" contains only one
// of the query's content words (meeting), while "Meeting ID 795 777
// 3585" contains two (meeting, id) - a clear, explainable, correct
// margin, not a fragile similarity score. This is deliberately a
// literal, deterministic word-overlap count, not a semantic
// technique - it will not catch paraphrases with no shared
// vocabulary, but everything it does match, it matches for a reason
// that can be shown and checked, unlike a cosine similarity score.
//
// When the winner isn't unique (a tie, or zero words matched
// anywhere), the caller should fall back to returning the whole
// chunk - the existing, safe, always-correct-if-verbose default.
// HasAnyContentWordOverlap reports whether at least one of the
// question's content words (stopwords stripped) appears anywhere in
// the combined content of the given chunks. This is the safety gate
// that stops /kb ask from confidently answering a question that has
// nothing to do with the knowledge base at all.
//
// Root cause this fixes: SemanticSearch always returns its top-ranked
// chunks, with no minimum similarity threshold - "most similar of
// what exists" is not the same as "actually relevant". A completely
// unrelated question (e.g. "what is the capital of Delhi" against a
// Sanskrit class schedule) was confirmed, in real testing, to return
// a confidently-formatted answer built from totally irrelevant
// content, with no indication anything was wrong.
//
// This does NOT use an embedding similarity cutoff. A prior
// investigation (see manager.go's history / CHANGELOG) found that
// this embedding model's cosine similarity scores don't reliably
// separate related from unrelated content even at the sentence
// level - an unrelated sentence scored HIGHER than a genuinely
// relevant one in real measured data, purely from incidental word
// overlap. Guessing a numeric threshold here, right before a
// release, would repeat that exact mistake. Content-word overlap is
// the same deterministic, explainable, already-proven-correct
// technique used by findPreciseLine and the "what is the" plural
// override - if a question shares literally zero words with
// anything retrieved, that's a safe, low-risk signal that nothing
// relevant was actually found, regardless of what embedding
// similarity ranked highest.
func HasAnyContentWordOverlap(question string, results []SemanticSearchResult) bool {
	words := contentWords(question)

	// Very short words (2 characters or fewer) are excluded from
	// this specific check: they collide as substrings inside
	// unrelated longer words far too easily to be a meaningful
	// relevance signal. Confirmed in real testing: "what is the
	// zoom meeting id" against a document containing no Zoom
	// information at all still passed this gate, because "id"
	// matched as a literal substring inside "Confidentiality" - a
	// coincidental match with zero relation to actual relevance.
	// findPreciseLine (a different function, used only to select a
	// single line within an already-relevant chunk, not to gate
	// relevance in the first place) is deliberately left unchanged -
	// its existing behavior is already verified correct for its
	// actual use case, and this length filter is scoped only to
	// the broader relevance-gating context where the false-positive
	// was found.
	var significantWords []string
	for _, w := range words {
		if len(w) > 2 {
			significantWords = append(significantWords, w)
		}
	}

	if len(significantWords) == 0 {
		// No question with meaningful content words at all can't be
		// gated - don't block it.
		return true
	}

	for _, result := range results {
		lower := strings.ToLower(result.Content)
		for _, w := range significantWords {
			if strings.Contains(lower, w) {
				return true
			}
		}
	}

	return false
}

func findPreciseLine(question, chunkContent string) (line string, found bool) {
	words := contentWords(question)
	if len(words) == 0 {
		return "", false
	}

	paragraphs := splitParagraphs(chunkContent)

	bestScore := 0
	bestIndex := -1
	tie := false

	for i, p := range paragraphs {
		lower := strings.ToLower(p)
		score := 0
		for _, w := range words {
			if strings.Contains(lower, w) {
				score++
			}
		}

		if score > bestScore {
			bestScore = score
			bestIndex = i
			tie = false
		} else if score == bestScore && score > 0 {
			tie = true
		}
	}

	if bestIndex == -1 || tie {
		return "", false
	}

	return paragraphs[bestIndex], true
}

// FilterRelevantChunks removes chunks that share zero content words
// with the question, operating at whole-chunk granularity only.
//
// This is deliberately NOT the same thing as filtering individual
// items within a chunk (see askEnumeration's doc comment in
// knowledge_handler.go for why that was tried and rejected: it
// silently dropped real, related entries - e.g. "Tuesday Chanting" -
// from a chunk about Sanskrit classes, just because the entry itself
// didn't contain the word "Sanskrit"). A chunk is included or
// excluded as a whole, so every item inside an included chunk is
// still returned together, exactly as before - this only decides
// whether an entire chunk belongs in the answer at all.
//
// Why this is needed: SemanticSearch always returns every chunk it
// has whenever the knowledge base's total chunk count is at or below
// the retrieval limit (5) - true for any KB with only 2-3 documents,
// regardless of embedding quality. HasAnyContentWordOverlap already
// gates the case where NOTHING in the results is relevant (a
// genuinely out-of-scope question), but it was never designed to
// remove individual irrelevant chunks from an otherwise-valid batch -
// confirmed in real testing on real Ollama embeddings: a question
// about a Sanskrit class schedule returned an entirely unrelated
// service contract alongside the correct answer, because both
// happened to be in the same small knowledge base.
//
// If filtering would remove every chunk, the original unfiltered
// list is returned instead - this should not happen in practice
// (HasAnyContentWordOverlap already confirms at least one chunk has
// overlap before this is ever called), but returning something
// verbose is always safer than returning nothing.
func FilterRelevantChunks(question string, results []SemanticSearchResult) []SemanticSearchResult {
	words := contentWords(question)

	var significantWords []string
	for _, w := range words {
		if len(w) > 2 {
			significantWords = append(significantWords, w)
		}
	}

	if len(significantWords) == 0 {
		return results
	}

	var filtered []SemanticSearchResult
	for _, result := range results {
		lower := strings.ToLower(result.Content)
		for _, w := range significantWords {
			if strings.Contains(lower, w) {
				filtered = append(filtered, result)
				break
			}
		}
	}

	if len(filtered) == 0 {
		return results
	}

	return filtered
}