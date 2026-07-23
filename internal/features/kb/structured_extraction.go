package kb

import (
	"fmt"
	"strings"
)

// ExtractedItem is one deterministically-extracted entry from a
// retrieved chunk, with source attribution preserved.
type ExtractedItem struct {
	Text             string
	SourceDocumentID string
	SourceChunkID    string
}

// ExtractItems returns each retrieved chunk as one item, with source
// attribution preserved.
//
// An earlier version of this function called splitParagraphs on each
// chunk's content, on the assumption that one blank-line-separated
// block corresponds to one logical entry (e.g. one class listing).
// That assumption was checked against a synthetic test fixture but
// never verified against a real document's actual formatting
// convention - and it was wrong. A real noticeboard-style .txt file
// was found, during v0.9.1 validation, to place a blank line between
// nearly every sentence (61 blank-line blocks for 9 real entries),
// not one blank line per entry. Calling splitParagraphs on chunk
// content therefore fragmented single entries into many meaningless
// sentence-level bullets ("The term begins on 12 January 2026." as
// its own bullet, unattached to which class it belongs to).
//
// Chunk boundaries are a safer unit to extract on: ChunkDocument
// already groups related content into ~DefaultChunkSize blocks
// without corrupting or splitting any single entry mid-content (see
// chunker_test.go). Treating each chunk as one item is coarser -
// multiple entries can share one chunk/bullet if they were grouped
// into the same chunk at import time - but it is provably correct:
// no entry is ever fragmented below the granularity chunking already
// guaranteed was safe.
func ExtractItems(results []SemanticSearchResult) []ExtractedItem {
	var items []ExtractedItem

	for _, result := range results {
		items = append(items, ExtractedItem{
			Text:             result.Content,
			SourceDocumentID: result.DocumentID,
			SourceChunkID:    result.ChunkID,
		})
	}

	return items
}

// BuildListAnswer formats extracted items as a deterministic,
// zero-hallucination bullet list with source citations - no LLM
// involved. Every item ExtractItems found is guaranteed to appear
// here; there is no generation step that could silently drop one.
func BuildListAnswer(items []ExtractedItem) string {
	if len(items) == 0 {
		return "No relevant knowledge found to answer this question."
	}

	var builder strings.Builder

	builder.WriteString(fmt.Sprintf(
		"Here is everything found in your knowledge base across %d relevant section(s):\n\n",
		len(items),
	))

	seenSources := make(map[string]bool)
	var orderedSources []string

	for _, item := range items {
		builder.WriteString("* ")
		// Item text can be multi-line (a paragraph with detail lines
		// beneath the main entry) - keep it intact rather than
		// collapsing to a single line, since detail lines (fees,
		// coach names, room numbers, etc.) are themselves matching
		// content a user asked to see, not noise to strip.
		builder.WriteString(strings.ReplaceAll(item.Text, "\n", "\n  "))
		builder.WriteString("\n\n")

		if !seenSources[item.SourceChunkID] {
			seenSources[item.SourceChunkID] = true
			orderedSources = append(orderedSources, item.SourceChunkID)
		}
	}

	builder.WriteString("Sources:\n")
	for _, src := range orderedSources {
		builder.WriteString("  [")
		builder.WriteString(src)
		builder.WriteString("]\n")
	}

	return builder.String()
}
