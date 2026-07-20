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

	builder.WriteString("Answer using only the supplied knowledge. If the question asks for multiple items (a list of timings, dates, entries, or similar), you must include every single matching item found across every source above - do not summarize, shorten, or omit any matching entry, even if the list is long. If the knowledge base does not contain enough information to answer, say so explicitly rather than guessing.")

	return builder.String()
}