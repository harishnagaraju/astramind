package kb

import (
	"strings"
	"testing"
)

func TestIsEnumerationQuery(t *testing.T) {
	cases := []struct {
		question string
		want     bool
	}{
		// Classic enumeration phrasings.
		{"What are all the Sanskrit class timings?", true},
		{"list out sanskrit classes timings", true},
		{"List every class or session mentioned", true},
		{"every single matching entry please", true},
		{"tell me all the fees", true},
		{"List the class timings", true},

		// The exact phrasing that slipped through the old, narrower
		// pattern list in real manual testing - no "all"/"every"/
		// "list" keyword at all, but still clearly asking for a list.
		// This is the regression case this rewrite exists to fix.
		{"what are the Sanskrit class timings", true},
		{"what are the club and class timings", true},

		// Clear single-fact lookups - the only cases that should
		// NOT be broadened.
		{"What time is Thursday Senior Sanskrit?", false},
		{"Is there a fee for the robotics workshop?", false},
		{"What is the zoom meeting ID?", false},
		{"Does the Friday class meet on holidays?", false},
		{"How much does the materials fee cost?", false},
		{"Who is the coach for chess club?", false},
		{"When is the next term starting?", false},
		{"", false},
	}

	for _, c := range cases {
		got := IsEnumerationQuery(c.question)
		if got != c.want {
			t.Errorf("IsEnumerationQuery(%q) = %v, want %v", c.question, got, c.want)
		}
	}
}

func TestExpandQuery_EnumerationQuestionIsExpanded(t *testing.T) {
	question := "What are all the Sanskrit class timings?"

	expanded := ExpandQuery(question)

	if expanded == question {
		t.Fatal("expected enumeration question to be expanded, got unchanged text")
	}

	if !strings.Contains(expanded, question) {
		t.Fatalf("expected expanded query to retain original question text, got %q", expanded)
	}

	if !strings.Contains(expanded, "related items") {
		t.Fatalf("expected expanded query to include broadening language, got %q", expanded)
	}
}

// TestExpandQuery_RegressionNoAllKeyword is the specific case found
// during manual testing: "what are the Sanskrit class timings" has no
// "all", "every", or "list" keyword, but is clearly asking for every
// matching class - and under the old narrow pattern list this was
// left unexpanded, silently dropping 2 of 9 real entries (Tuesday
// Chanting, Tuesday Gita Youth group) in two consecutive real runs
// through /kb ask.
func TestExpandQuery_RegressionNoAllKeyword(t *testing.T) {
	question := "what are the Sanskrit class timings"

	expanded := ExpandQuery(question)

	if expanded == question {
		t.Fatal("expected this question to be expanded (regression: this exact phrasing was previously missed)")
	}
}

func TestExpandQuery_SingleFactQuestionIsUnchanged(t *testing.T) {
	question := "What time is Thursday Senior Sanskrit?"

	expanded := ExpandQuery(question)

	if expanded != question {
		t.Fatalf("expected single-fact question to be returned unchanged, got %q", expanded)
	}
}

func TestExpandQuery_EmptyQuestionIsUnchanged(t *testing.T) {
	if got := ExpandQuery(""); got != "" {
		t.Fatalf("expected empty question to be returned unchanged, got %q", got)
	}

	if got := ExpandQuery("   "); got != "   " {
		t.Fatalf("expected whitespace-only question to be returned unchanged, got %q", got)
	}
}
