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

// TestChunkRespectsParagraphBoundaries guards against the class of bug
// that originally caused "Tuesday Gita Youth group" to come back as
// "uth group" in a real /kb ask answer: the old byte-offset splitter
// had no awareness of word or paragraph boundaries and would slice
// straight through the middle of a word whenever a cut point happened
// to land there.
//
// This test builds a document out of distinct, known blocks separated
// by blank lines (the same structure as a real noticeboard/schedule
// document) and asserts two things mechanically, with no LLM involved:
//
//  1. Every known identifying phrase appears completely intact in at
//     least one chunk.
//  2. No chunk starts mid-word - i.e. at a position that splits a word
//     which existed as one unbroken token in the source content.
//
// This is the cheapest possible layer to catch this bug: no server,
// no model, no manual script run - just `go test`.
func TestChunkRespectsParagraphBoundaries(t *testing.T) {
	blocks := []string{
		"Monday Chess Club 15:00 - 16:00.\nHeld in Room 4. Open to all students grade 6 and up.\nCoach: Mr. Fernandes. No fee required.",
		"Tuesday Robotics Workshop 14:00 - 15:30.\nHeld in the Science Lab. Registration required in advance.\nCoach: Ms. Alvares. Materials fee: 10 dollars.",
		"Wednesday Debate Team 16:00 - 17:00.\nHeld in the Auditorium. Open to grade 8 and up only.\nCoach: Mr. Pinto. No fee required.",
		"Thursday Art Class 13:00 - 14:00.\nHeld in Room 9. Open to all students.\nCoach: Ms. Rodrigues. Materials fee: 5 dollars.",
		"Friday Music Rehearsal 15:30 - 16:30.\nHeld in the Music Room. Open to band members only.\nCoach: Mr. Souza. No fee required.",
		"Saturday Yoga Session 09:00 - 10:00.\nHeld in the Gym. Open to all staff and students.\nCoach: Ms. Costa. No fee required.",
		"Sunday Study Hall 10:00 - 12:00.\nHeld in the Library. Open to all students preparing for exams.\nCoach: Mr. Dias. No fee required.",
		"Monday Evening Coding Club 17:30 - 18:30.\nHeld in Room 12, the computer lab. Open to grade 7 and up.\nCoach: Mr. Lobo. No fee required.",
	}

	content := strings.Join(blocks, "\n\n")

	// Content is deliberately long enough (with these blocks) to
	// exceed a small chunk size, forcing the splitter to actually
	// make a cut somewhere - a document short enough to fit in one
	// chunk can never exercise this bug, which is exactly why the
	// original single-sentence fixtures (see TestChunkSmallDocument)
	// never caught it.
	const chunkSize = 300
	const overlap = 50

	doc := &Document{ID: "doc1", Content: content}
	chunks := ChunkDocument(doc, chunkSize, overlap)

	if len(chunks) < 2 {
		t.Fatalf("expected content to be split into multiple chunks, got %d", len(chunks))
	}

	identifiers := []string{
		"Chess Club",
		"Robotics Workshop",
		"Debate Team",
		"Art Class",
		"Music Rehearsal",
		"Yoga Session",
		"Study Hall",
		"Coding Club",
	}

	for _, id := range identifiers {
		found := false
		for _, c := range chunks {
			if strings.Contains(c.Content, id) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("identifier %q not found intact in any chunk - likely split mid-word", id)
		}
	}

	for i, c := range chunks {
		if len(c.Content) == 0 {
			t.Errorf("chunk %d is empty", i)
			continue
		}

		idx := strings.Index(content, c.Content)
		if idx <= 0 {
			continue
		}

		first := c.Content[0]
		prev := content[idx-1]

		isLetter := func(b byte) bool {
			return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z')
		}

		if first >= 'a' && first <= 'z' && isLetter(prev) {
			end := 20
			if end > len(c.Content) {
				end = len(c.Content)
			}
			t.Errorf("chunk %d appears to start mid-word: %q...", i, c.Content[:end])
		}
	}
}

// TestChunkHandlesCRLFLineEndings guards against a bug found during
// v0.9.1 validation: real-world files (especially anything authored
// or saved on Windows) commonly use CRLF ("\r\n") line endings. A
// CRLF blank line is "\r\n\r\n", which never contains the substring
// "\n\n" - so paragraph splitting would silently find zero boundaries
// in a CRLF file, collapse the whole document into one oversized
// "paragraph", and fall through to the old byte-offset hardSplit
// fallback, reproducing the exact mid-word corruption bug this
// package's chunker was supposed to have already fixed.
//
// This is exactly the class of gap TestChunkRespectsParagraphBoundaries
// could not catch: that test's fixture is built from Go string
// literals joined with "\n", which is LF, not CRLF - so it never
// exercised this path. A real user-supplied .txt file with CRLF
// endings slipped through every existing automated test while still
// being broken, which is why this test exists as its own case rather
// than assuming CRLF is covered by the LF version.
func TestChunkHandlesCRLFLineEndings(t *testing.T) {
	blocks := []string{
		"Monday Chess Club 15:00 - 16:00.\r\nHeld in Room 4. Open to all students grade 6 and up.\r\nCoach: Mr. Fernandes. No fee required.",
		"Tuesday Robotics Workshop 14:00 - 15:30.\r\nHeld in the Science Lab. Registration required in advance.\r\nCoach: Ms. Alvares. Materials fee: 10 dollars.",
		"Wednesday Debate Team 16:00 - 17:00.\r\nHeld in the Auditorium. Open to grade 8 and up only.\r\nCoach: Mr. Pinto. No fee required.",
		"Thursday Art Class 13:00 - 14:00.\r\nHeld in Room 9. Open to all students.\r\nCoach: Ms. Rodrigues. Materials fee: 5 dollars.",
		"Friday Music Rehearsal 15:30 - 16:30.\r\nHeld in the Music Room. Open to band members only.\r\nCoach: Mr. Souza. No fee required.",
		"Saturday Yoga Session 09:00 - 10:00.\r\nHeld in the Gym. Open to all staff and students.\r\nCoach: Ms. Costa. No fee required.",
		"Sunday Study Hall 10:00 - 12:00.\r\nHeld in the Library. Open to all students preparing for exams.\r\nCoach: Mr. Dias. No fee required.",
		"Monday Evening Coding Club 17:30 - 18:30.\r\nHeld in Room 12, the computer lab. Open to grade 7 and up.\r\nCoach: Mr. Lobo. No fee required.",
	}

	// CRLF blank-line separator between blocks, matching what a real
	// Windows-authored .txt file looks like on disk.
	content := strings.Join(blocks, "\r\n\r\n")

	const chunkSize = 300
	const overlap = 50

	doc := &Document{ID: "doc1", Content: content}
	chunks := ChunkDocument(doc, chunkSize, overlap)

	if len(chunks) < 2 {
		t.Fatalf("expected content to be split into multiple chunks, got %d", len(chunks))
	}

	identifiers := []string{
		"Chess Club",
		"Robotics Workshop",
		"Debate Team",
		"Art Class",
		"Music Rehearsal",
		"Yoga Session",
		"Study Hall",
		"Coding Club",
	}

	for _, id := range identifiers {
		found := false
		for _, c := range chunks {
			if strings.Contains(c.Content, id) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("identifier %q not found intact in any chunk on CRLF input - likely split mid-word", id)
		}
	}
}
