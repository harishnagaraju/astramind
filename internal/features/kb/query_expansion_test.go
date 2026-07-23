package kb

import (
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

		// The exact phrasing that slipped through an earlier, narrower
		// pattern list in real manual testing - no "all"/"every"/
		// "list" keyword at all, but still clearly asking for a list.
		// This is the regression case the current default exists to
		// guard against.
		{"what are the Sanskrit class timings", true},
		{"what are the club and class timings", true},

		// Clear single-fact lookups - the only cases that should
		// route to single-fact extraction instead of enumeration.
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