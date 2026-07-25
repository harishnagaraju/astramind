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
