package kb

import (
	"regexp"
	"strings"
)

// singleFactPatterns matches question phrasings that are clearly
// asking for one specific fact about one specific named thing, as
// opposed to enumerating everything in a category. Only questions
// matching one of these are treated as single-fact lookups by
// IsEnumerationQuery; everything else defaults to enumeration.
//
// This list is deliberately narrow and deliberately the exception,
// not the rule - see IsEnumerationQuery's doc comment for why the
// default is "enumeration unless clearly single-fact" rather than
// "enumeration only if clearly asking for a list".
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

// IsEnumerationQuery reports whether a question should be treated as
// asking for a complete list/enumeration of matching items, as
// opposed to a single specific fact lookup. Used by handleKBAsk to
// route between the deterministic enumeration path (ExtractItems /
// BuildListAnswer) and the deterministic single-fact path
// (ExtractiveAnswer).
//
// The default is enumeration: a question is only treated as a
// single-fact lookup if it clearly matches one of singleFactPatterns
// (e.g. "is there a fee for X", "what time is X", "how much does X
// cost"). Everything else - including phrasings with no explicit
// "all"/"every"/"list" keyword at all, like "what are the Sanskrit
// class timings" - is treated as enumeration.
//
// This default was flipped deliberately after testing showed that an
// earlier, narrower version (which only matched explicit enumeration
// keywords like "what are all", "list every") missed extremely
// natural phrasings such as "what are the X timings" - which has no
// "all" in it at all, but is asking for a list just as much as "what
// are all the X timings" is. Given that the failure mode being
// guarded against is a silent, undetectable dropped fact, a false
// positive here (a single-fact question routed to the broader
// enumeration path) is a far smaller cost than a false negative (an
// enumeration question routed to single-fact extraction, silently
// dropping real entries).
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