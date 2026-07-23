package kb

import (
	"strings"
	"testing"
)

// realSanskritNoticeboard is the exact byte content of the real
// Sanskrit1.txt file used throughout v0.9.1 validation testing -
// CRLF line endings included, unmodified. This is not a synthetic
// fixture: it is the actual document that repeatedly exposed the
// chunking and enumeration bugs found during manual testing. Using
// it directly here means this test proves the fix against the real
// failure case, not a reconstruction of it.
const realSanskritNoticeboard = "NOTICEBOARD\r\n\r\n\r\n\r\nRutger on zoom. Most courses are 8 weeks.\r\n\r\nAll courses are on the same zoom address.\r\n\r\nLocal Dublin times are the same as GMT.\r\n\r\nhttps://us04web.zoom.us/j/7957773585?pwd=eExsNld3bXJxdzJ5ZFBKYUg5NHliZz09\r\n\r\nMeeting ID 795 777 3585\r\n\r\nPassword OMpeace\r\n\r\n\r\n\r\nMonday Sanskrit Term 14 Youth from 17:00-18:00.\r\n\r\nThe term begins on 12 January 2026.\r\n\r\nNot meeting on 16 February. Last week 9 March.\r\n\r\nNo fee. Open to anyone with enough knowledge,\r\n\r\nbut if you are new, write to rutger@johnscottus.ie\r\n\r\n[Enrol with VHCCI https://www.hindu.ie/]\r\n\r\n\r\n\r\nTuesday Chanting 06:30 - 07:15.\r\n\r\nThe next term begins on 13 January 2026.\r\n\r\nNot meeting on 17 February. Last week 10 March.\r\n\r\nNo fee. Open to anyone, but if you are new, write to rutger@johnscottus.ie\r\n\r\n\r\n\r\nTuesday G\u012bt\u0101 Youth group 17:00 - 18:00.\r\n\r\nThe next term begins on 13 January 2026.\r\n\r\nNot meeting on 17 February. Last week 10 March.\r\n\r\n[Register with VHCCI https://www.hindu.ie/]\r\n\r\nOpen to young people. No fee. Contact rutger@johnscottus.ie\r\n\r\n\r\n\r\nWednesday Chanting 06:30 - 07:15.\r\n\r\nThe next term begins on 14 January 2026.\r\n\r\nNot meeting on 18 February. Last week 11 March.\r\n\r\nNo fee. Open to anyone, but if you are new, contact rutger@johnscottus.ie\r\n\r\n\r\n\r\nWednesday Sanskrit Term 4 from 17:00-18:00.\r\n\r\nThe term begins on 14 January 2026.\r\n\r\nNot meeting on 18 February. Last week 11 March.\r\n\r\nChildren registration portal with VHCCI https://www.hindu.ie/\r\n\r\nAdult registration link https://practicalphilosophy.ie/enrol/sol/sk-in12.\r\n\r\n\r\n\r\nWednesday G\u012bt\u0101 group Evenings 19:00-20:00.\r\n\r\nThe next term begins on 14 January 2025.\r\n\r\nNot meeting on 18 February. Last week 11 March.\r\n\r\nOnly open to SPES. No fee. For information, contact rutger@johnscottus.ie\r\n\r\n\r\n\r\nThursday Senior Sanskrit from 06:15 to 07:15 am.\r\n\r\nThe next term begins on 15 January 2026.\r\n\r\nNot meeting on 19 February. Last week 12 March.\r\n\r\nNo fee for SPES students. \u20ac 60 to enrol on https://practicalphilosophy.ie/enrol/sol/sk-in12\r\n\r\nOpen to anyone, but if new, contact rutger@johnscottus.ie\r\n\r\n\r\n\r\nFriday Term 4 Sanskrit from 18:30 to 19:30.\r\n\r\nThe term will start on 16 January 2026.\r\n\r\nNot meeting on 20 February. Last week 13 March.\r\n\r\nRun by the Eire Vedanta Society. Enrol using https://www.rkmireland.org/\r\n\r\n\r\n\r\nSaturday Sanskrit Term 17 from 17:00 - 18:00.\r\n\r\nThe term begins on 17 January 2026.\r\n\r\nNot meeting on 21 February. Last week 14 March.\r\n\r\nNo fee for SPES students. \u20ac 60 to enrol on https://practicalphilosophy.ie/enrol/sol/sk-in12\r\n\r\nOpen to anyone. rutger@johnscottus.ie\r\n\r\n\r\n\r\nLeaf\r\nWELL BEING\r\nSVASTI\r\n"

// TestExtractItems_RealSanskritDocument is the deterministic
// replacement for the manual "run /kb ask 5 times and eyeball the
// output" acceptance test that was used throughout this
// investigation. It runs in milliseconds, involves no LLM, and
// either passes or fails with no ambiguity - unlike the LLM-based
// path, there is no run-to-run variance to interpret.
//
// It proves, end to end through the real pipeline (chunking ->
// retrieval simulation -> extraction -> formatting), that all 9 real
// class entries are recovered, regardless of how the question is
// worded - because the deterministic path never filters items by
// keyword overlap with the question at all.
func TestExtractItems_RealSanskritDocument(t *testing.T) {
	doc := &Document{ID: "sanskrit-doc", Content: realSanskritNoticeboard}
	chunks := ChunkDocument(doc, DefaultChunkSize, DefaultOverlap)

	if len(chunks) == 0 {
		t.Fatal("expected at least one chunk from the real document")
	}

	// Simulate what SemanticSearch would return: since this document
	// produces fewer chunks than DefaultSemanticSearchLimit, every
	// chunk is always retrieved regardless of query wording - this
	// mirrors the real retrieval behavior exactly, not an idealized
	// version of it.
	var results []SemanticSearchResult
	for _, c := range chunks {
		results = append(results, SemanticSearchResult{
			DocumentID: c.DocumentID,
			ChunkID:    c.ID,
			ChunkIndex: c.Index,
			Content:    c.Content,
		})
	}

	items := ExtractItems(results)

	// Strict count check: this is what the earlier version of this
	// test failed to assert, and why it passed despite ExtractItems
	// being broken. Checking "the right substring is present
	// somewhere" is not sufficient to catch over-fragmentation - a
	// substring match succeeds whether item.Text is one coherent
	// entry or one meaningless sentence fragment. The count must
	// equal the number of chunks, not the number of sentences.
	if len(items) != len(chunks) {
		t.Fatalf(
			"expected exactly %d items (one per chunk), got %d - "+
				"chunk content is being fragmented below chunk granularity",
			len(chunks), len(items),
		)
	}

	// The 9 real, distinct class/session entries in the actual file.
	// This list was built by reading the real document directly, not
	// copied from a prior (possibly incomplete) LLM answer.
	expectedEntries := []string{
		"Monday Sanskrit Term 14 Youth from 17:00-18:00",
		"Tuesday Chanting 06:30 - 07:15",
		"Tuesday Gītā Youth group 17:00 - 18:00",
		"Wednesday Chanting 06:30 - 07:15",
		"Wednesday Sanskrit Term 4 from 17:00-18:00",
		"Wednesday Gītā group Evenings 19:00-20:00",
		"Thursday Senior Sanskrit from 06:15 to 07:15 am",
		"Friday Term 4 Sanskrit from 18:30 to 19:30",
		"Saturday Sanskrit Term 17 from 17:00 - 18:00",
	}

	for _, expected := range expectedEntries {
		found := false
		for _, item := range items {
			if strings.Contains(item.Text, expected) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("entry not found in extracted items: %q", expected)
		}
	}

	// Also confirm the formatted answer contains every entry - this
	// is what a real /kb ask response would show the user.
	answer := BuildListAnswer(items)
	for _, expected := range expectedEntries {
		if !strings.Contains(answer, expected) {
			t.Errorf("entry missing from formatted answer: %q", expected)
		}
	}
}

func TestBuildListAnswer_EmptyItems(t *testing.T) {
	answer := BuildListAnswer(nil)
	if !strings.Contains(answer, "No relevant knowledge found") {
		t.Fatalf("expected empty-items message, got %q", answer)
	}
}

func TestBuildListAnswer_IncludesSourceCitations(t *testing.T) {
	items := []ExtractedItem{
		{Text: "Entry one", SourceDocumentID: "doc1", SourceChunkID: "chunk-a"},
		{Text: "Entry two", SourceDocumentID: "doc1", SourceChunkID: "chunk-b"},
	}

	answer := BuildListAnswer(items)

	if !strings.Contains(answer, "Entry one") || !strings.Contains(answer, "Entry two") {
		t.Fatal("expected both entries in formatted answer")
	}

	if !strings.Contains(answer, "[chunk-a]") || !strings.Contains(answer, "[chunk-b]") {
		t.Fatal("expected both source chunk IDs cited in formatted answer")
	}
}
