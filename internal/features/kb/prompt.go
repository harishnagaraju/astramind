package kb

import (
	"strconv"
	"strings"
)

// BuildPrompt builds a RAG prompt from a user question and search results.
func BuildPrompt(question string, results []SearchResult) string {
	var builder strings.Builder

	builder.WriteString("You are answering questions using the provided knowledge base.\n\n")

	builder.WriteString("Knowledge Base:\n\n")

	for _, result := range results {
		builder.WriteString("[Document: ")
		builder.WriteString(result.DocumentID)
		builder.WriteString("]\n")

		builder.WriteString(result.Content)
		builder.WriteString("\n\n")
	}

	builder.WriteString("Question:\n")
	builder.WriteString(question)
	builder.WriteString("\n\n")

	builder.WriteString("Answer using only the supplied knowledge.")

	return builder.String()
}

// BuildSemanticPrompt builds a RAG prompt from a user question and
// semantic (embedding-based) search results. Mirrors BuildPrompt's
// format exactly, so downstream rendering doesn't need to know which
// search mode produced the results.
//
// The prompt explicitly numbers the sources and instructs the model
// to address every one of them in turn, rather than relying on a
// single trailing "include everything" instruction. This matters
// because a trailing instruction is a comparatively weak signal
// competing against the literal wording of the question - testing
// during v0.9.1 validation showed a model would consistently follow
// the question's exact keywords (e.g. answer "Sanskrit class"
// narrowly) even with a broadening instruction present in the same
// prompt, especially at low temperature where the model reliably
// follows its single strongest signal rather than exploring weaker
// ones. Requiring the model to work through each numbered source
// individually is a stronger, structural way to prevent silent
// source omission than instruction wording alone.
func BuildSemanticPrompt(question string, results []SemanticSearchResult) string {
	var builder strings.Builder

	builder.WriteString("You are answering questions using the provided knowledge base.\n\n")

	builder.WriteString("Knowledge Base:\n\n")

	for i, result := range results {
		builder.WriteString("[Source ")
		builder.WriteString(strconv.Itoa(i + 1))
		builder.WriteString(" of ")
		builder.WriteString(strconv.Itoa(len(results)))
		builder.WriteString(", Document: ")
		builder.WriteString(result.DocumentID)
		builder.WriteString("]\n")

		builder.WriteString(result.Content)
		builder.WriteString("\n\n")
	}

	builder.WriteString("Question:\n")
	builder.WriteString(question)
	builder.WriteString("\n\n")

	builder.WriteString("Answer using only the supplied knowledge. ")
	builder.WriteString("There are ")
	builder.WriteString(strconv.Itoa(len(results)))
	builder.WriteString(" sources above, numbered 1 through ")
	builder.WriteString(strconv.Itoa(len(results)))
	builder.WriteString(". ")
	builder.WriteString("Work through each source in order and check it individually for any information relevant to the question, even if that source doesn't use the same words as the question - a source about a related topic, activity, or category is still relevant even without an exact keyword match. ")
	builder.WriteString("If the question asks for multiple items (a list of timings, dates, entries, or similar), you must include every single matching item found in every source, one by one - do not summarize, shorten, or skip a source, even if the list is long. ")
	builder.WriteString("If the knowledge base does not contain enough information to answer, say so explicitly rather than guessing.")

	return builder.String()
}
