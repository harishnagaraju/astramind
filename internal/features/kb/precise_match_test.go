package kb

import "testing"

// TestFindPreciseLine_ZoomMeetingID reproduces the exact real case
// that defeated embedding-similarity-based windowing earlier in this
// project's history: "Not meeting on 16 February" scored HIGHER
// cosine similarity against "what is the zoom meeting id" than the
// genuinely relevant "Meeting ID 795 777 3585" did, because both
// happened to contain the word "meeting" and the embedding model had
// no reliable way to weigh that against real relevance. Content-word
// overlap counting does not have this failure mode: the real answer
// contains two of the question's content words (meeting, id) while
// the false-positive contains only one (meeting) - a clear,
// explainable margin.
func TestFindPreciseLine_ZoomMeetingID(t *testing.T) {
	content := `NOTICEBOARD

Rutger on zoom. Most courses are 8 weeks.

All courses are on the same zoom address.

https://us04web.zoom.us/j/7957773585

Meeting ID 795 777 3585

Password OMpeace

Monday Sanskrit Term 14 Youth from 17:00-18:00.

The term begins on 12 January 2026.

Not meeting on 16 February. Last week 9 March.`

	line, found := findPreciseLine("what is the zoom meeting id", content)
	if !found {
		t.Fatal("expected a unique precise match")
	}
	if line != "Meeting ID 795 777 3585" {
		t.Fatalf("expected exact meeting ID line, got %q", line)
	}
}

// TestFindPreciseLine_ThursdaySanskrit reproduces the other real
// case verified earlier: correctly distinguishing Thursday's entry
// from Wednesday's and Friday's, which all share the word "sanskrit".
func TestFindPreciseLine_ThursdaySanskrit(t *testing.T) {
	content := `Wednesday Sanskrit Term 4 from 17:00-18:00.

Thursday Senior Sanskrit from 06:15 to 07:15 am.

Friday Term 4 Sanskrit from 18:30 to 19:30.`

	line, found := findPreciseLine("what time is thursday senior sanskrit", content)
	if !found {
		t.Fatal("expected a unique precise match")
	}
	if line != "Thursday Senior Sanskrit from 06:15 to 07:15 am." {
		t.Fatalf("expected the Thursday line, got %q", line)
	}
}

// TestFindPreciseLine_TieFallsBack confirms that when two paragraphs
// score equally, findPreciseLine reports no match at all - the
// caller is expected to fall back to the safe whole-chunk return
// rather than arbitrarily pick one of the tied candidates.
func TestFindPreciseLine_TieFallsBack(t *testing.T) {
	content := `Tuesday Chanting 06:30 - 07:15.

Wednesday Chanting 06:30 - 07:15.`

	// Both paragraphs contain "chanting" and nothing else
	// distinguishing - a genuine tie.
	_, found := findPreciseLine("what time is chanting", content)
	if found {
		t.Fatal("expected no unique match on a tie, got one anyway")
	}
}

// TestFindPreciseLine_NoMatchFallsBack confirms zero content-word
// overlap anywhere also correctly reports no match.
func TestFindPreciseLine_NoMatchFallsBack(t *testing.T) {
	content := `Monday Chess Club 15:00 - 16:00.

Tuesday Robotics Workshop 14:00 - 15:30.`

	_, found := findPreciseLine("what is the capital of France", content)
	if found {
		t.Fatal("expected no match for an unrelated question, got one anyway")
	}
}

func TestFindPreciseLine_EmptyQuestion(t *testing.T) {
	_, found := findPreciseLine("", "Some content here.")
	if found {
		t.Fatal("expected no match for an empty question")
	}
}

func TestContentWords_StripsStopwordsAndPunctuation(t *testing.T) {
	words := contentWords("What is the zoom meeting ID?")
	expected := []string{"zoom", "meeting", "id"}

	if len(words) != len(expected) {
		t.Fatalf("expected %v, got %v", expected, words)
	}
	for i, w := range expected {
		if words[i] != w {
			t.Fatalf("expected %v, got %v", expected, words)
		}
	}
}
