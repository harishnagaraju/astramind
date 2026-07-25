package kb

import (
	"regexp"
	"strings"
)

// whatIsThePattern is checked separately from the other
// singleFactPatterns, because it is the one pattern in the list that
// is genuinely ambiguous between "single fact" and "enumeration"
// depending on what follows it. "what is the zoom meeting id" is a
// single fact. "what is the timings of sanskrit classes" is asking
// for every timing, despite starting with the identical three words
// - a real bug found in manual testing, where this phrasing was
// misrouted to single-fact extraction and silently dropped 5 of 9
// real entries because ExtractiveAnswer only ever returns one chunk.
var whatIsThePattern = regexp.MustCompile(`(?i)^\s*what is the\s+(\S+)`)

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
	whatIsThePattern,
	regexp.MustCompile(`(?i)^\s*how much\b`),
	regexp.MustCompile(`(?i)^\s*who is\b`),
	regexp.MustCompile(`(?i)^\s*where is\b`),
}

// looksPlural is a deliberately simple heuristic (ends in "s", not
// "ss") for whether a word is likely a plural noun. It is not
// linguistically robust - words like "robotics" or "physics" would
// be misclassified - but it is only used to decide whether to widen
// an already-ambiguous match to the safer (enumeration) path, never
// to narrow one. A false positive here (treating a genuinely
// singular question as enumeration) costs a slightly more verbose
// answer; a false negative (missing a real plural) costs a silently
// incomplete answer, which is the actual bug this exists to prevent.
func looksPlural(word string) bool {
	lower := strings.ToLower(word)
	return strings.HasSuffix(lower, "s") && !strings.HasSuffix(lower, "ss")
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
		if !p.MatchString(trimmed) {
			continue
		}

		// "what is the" is ambiguous - only treat it as single-fact
		// if the word right after "the" doesn't look plural. Every
		// other pattern in the list (is there, how much, who is,
		// etc.) carries a clearer single-fact signal on its own and
		// isn't second-guessed here.
		if p == whatIsThePattern {
			match := whatIsThePattern.FindStringSubmatch(trimmed)
			if len(match) == 2 && looksPlural(match[1]) {
				continue
			}
		}

		return false
	}

	return true
}
