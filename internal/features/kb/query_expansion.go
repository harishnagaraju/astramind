package kb

import (
	"regexp"
	"strings"
)

// enumerationPatterns matches question phrasings that ask for a
// complete list/enumeration of matching items, as opposed to a
// single-fact lookup. These are the questions most exposed to the
// narrow-keyword-anchoring bug: when the question echoes only some
// KB entries' exact wording (e.g. "Sanskrit class timings"), a
// smaller model tends to retrieve and report only the entries whose
// text closely matches the question, silently omitting related but
// differently-worded entries (e.g. "Chanting", "Gita group" sessions
// that are Sanskrit-related but don't contain the word "class").
//
// Confirmed via direct model testing (bypassing AstraMind entirely):
// narrow phrasing "what are the Sanskrit class timings" recovered
// 4 of 9 real entries, reproducibly, across 5 runs. Broadening the
// question to "list every class or session" recovered 9 of 9,
// reproducibly, across 5 runs. A generic "interpret broadly" system
// prompt instruction alone was NOT sufficient - it only partially
// and inconsistently recovered borderline cases. Rewriting the
// question itself, before it reaches the model, is what worked.

// singleFactPatterns matches question phrasings that are clearly
// asking for one specific fact about one specific named thing, as
// opposed to enumerating everything in a category. These are the
// only questions ExpandQuery leaves untouched.
//
// This list is deliberately narrow and deliberately the exception,
// not the rule - see the package comment on ExpandQuery for why the
// default was flipped to "broaden unless clearly single-fact" rather
// than "broaden only if clearly asking for a list".
var singleFactPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)^\s*is there\b`),
	regexp.MustCompile(`(?i)^\s*does\b`),
	regexp.MustCompile(`(?i)^\s*is\b`),
	regexp.MustCompile(`(?i)^\s*was\b`),
	regexp.MustCompile(`(?i)^\s*when (is|does|will)\b`),
	regexp.MustCompile(`(?i)^\s*what time is\b`),
	regexp.MustCompile(`(?i)^\s*what is the\b`),
	regexp.MustCompile(`(?i)^\s*how much\b`),
	regexp.MustCompile(`(?i)^\s*who is\b`),
	regexp.MustCompile(`(?i)^\s*where is\b`),
}

// expansionSuffix is appended to the question text itself (not just
// the system prompt) for any question that isn't a clear single-fact
// lookup. Testing showed that broadening the QUESTION is what changes
// retrieval and generation behavior - a broadening instruction placed
// only in the system prompt was not enough on its own.
const expansionSuffix = ", including all related items, sessions, or entries that reasonably fall under this category, even if they don't use the exact same wording as this question"

// IsEnumerationQuery reports whether a question should be treated as
// asking for a complete list/enumeration of matching items, as
// opposed to a single specific fact lookup.
//
// The default is enumeration: a question is only treated as a
// single-fact lookup if it clearly matches one of singleFactPatterns
// (e.g. "is there a fee for X", "what time is X", "how much does X
// cost"). Everything else - including phrasings with no explicit
// "all"/"every"/"list" keyword at all, like "what are the Sanskrit
// class timings" - is treated as enumeration.
//
// This default was flipped deliberately after testing showed that an
// earlier, narrower version (which only expanded questions matching
// explicit enumeration keywords like "what are all", "list every")
// missed extremely natural phrasings such as "what are the X
// timings" - which has no "all" in it at all, but is asking for a
// list just as much as "what are all the X timings" is. Given that
// the failure mode being guarded against is a silent, undetectable
// dropped fact (see chunker.go / this package's test fixtures for the
// real-world case this was built from), a false positive here (a
// single-fact question gets slightly broadened) is a far smaller cost
// than a false negative (an enumeration question doesn't get
// broadened and silently drops real entries).
func IsEnumerationQuery(question string) bool {
	trimmed := strings.TrimSpace(question)
	if trimmed == "" {
		return false
	}

	for _, p := range singleFactPatterns {
		if p.MatchString(trimmed) {
			return false
		}
	}

	return true
}

// ExpandQuery rewrites questions to reduce narrow-keyword-anchoring
// bias before the question is used for semantic search and prompt
// construction. Only questions matching singleFactPatterns are left
// unchanged - everything else is broadened, on the assumption that
// most KB questions ask about a category ("Sanskrit classes", "fee
// schedule", "club timings") rather than a single named fact, and
// that broadening a genuinely single-fact question costs far less
// than silently dropping real entries from an enumeration question.
//
// This is a rule-based approach. If the single-fact pattern list
// proves too narrow (broadening genuinely hurts single-fact answer
// precision in practice), the next step is an LLM-based query
// classification/rewrite (a small, separate, low-temperature call)
// instead of maintaining this pattern list indefinitely.
func ExpandQuery(question string) string {
	trimmed := strings.TrimSpace(question)
	if trimmed == "" {
		return question
	}

	if !IsEnumerationQuery(trimmed) {
		return question
	}

	return trimmed + expansionSuffix
}
