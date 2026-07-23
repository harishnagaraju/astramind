package kb

import "strings"

const (
	DefaultChunkSize = 1000
	DefaultOverlap   = 200
)

// ChunkDocument splits doc.Content into chunks, preferring to break on
// blank-line (paragraph) boundaries so a chunk never cuts through the
// middle of a word or a logical block of text. If a single paragraph is
// itself larger than chunkSize, that paragraph falls back to the
// original byte-offset sliding-window split with overlap.
//
// This fixes a real bug found during v0.9.1 hardware/RAG validation:
// the original implementation sliced content at raw byte offsets with
// no awareness of word boundaries, producing corrupted mid-word chunks
// on real documents (e.g. "Tuesday Gita Youth group" -> "uth group")
// and, downstream, non-deterministic omissions in /kb ask answers.
func ChunkDocument(doc *Document, chunkSize, overlap int) []Chunk {
	if chunkSize <= 0 {
		chunkSize = DefaultChunkSize
	}

	if overlap < 0 {
		overlap = 0
	}

	// Normalize line endings before splitting. Real-world text files
	// - especially anything authored or saved on Windows - commonly
	// use CRLF ("\r\n") line endings. A CRLF blank line is "\r\n\r\n",
	// which never contains the substring "\n\n", so splitParagraphs
	// would silently find zero paragraph boundaries in such a file,
	// collapse the whole document into one oversized "paragraph", and
	// fall through to the byte-offset hardSplit fallback - reproducing
	// the exact mid-word corruption bug this function exists to fix.
	// Confirmed against a real CRLF-encoded fixture during v0.9.1
	// validation testing.
	normalized := strings.ReplaceAll(doc.Content, "\r\n", "\n")
	normalized = strings.ReplaceAll(normalized, "\r", "\n")

	paragraphs := splitParagraphs(normalized)

	var chunks []Chunk
	var current strings.Builder
	index := 0

	flush := func() {
		if current.Len() == 0 {
			return
		}
		chunks = append(chunks, Chunk{
			ID:         generateDocumentID(),
			DocumentID: doc.ID,
			Index:      index,
			Content:    current.String(),
		})
		index++
		current.Reset()
	}

	for _, para := range paragraphs {
		if len(para) > chunkSize {
			// Paragraph itself is too big - flush what we have, then
			// fall back to the original byte-offset sliding window
			// for this paragraph only.
			flush()
			for _, sub := range hardSplit(para, chunkSize, overlap) {
				chunks = append(chunks, Chunk{
					ID:         generateDocumentID(),
					DocumentID: doc.ID,
					Index:      index,
					Content:    sub,
				})
				index++
			}
			continue
		}

		// Would adding this paragraph overflow the current chunk?
		if current.Len() > 0 && current.Len()+len(para)+2 > chunkSize {
			flush()
		}

		if current.Len() > 0 {
			current.WriteString("\n\n")
		}
		current.WriteString(para)
	}

	flush()

	return chunks
}

// splitParagraphs splits content on blank lines, trims whitespace,
// and drops any empty entries.
func splitParagraphs(content string) []string {
	raw := strings.Split(content, "\n\n")

	var out []string
	for _, p := range raw {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

// hardSplit is the original fixed-size sliding-window splitter, kept
// as a fallback for paragraphs that individually exceed chunkSize.
func hardSplit(content string, chunkSize, overlap int) []string {
	var parts []string

	start := 0
	for start < len(content) {
		end := start + chunkSize
		if end > len(content) {
			end = len(content)
		}

		parts = append(parts, content[start:end])

		if end == len(content) {
			break
		}
		start = end - overlap
	}

	return parts
}
