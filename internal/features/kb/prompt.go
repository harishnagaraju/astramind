package kb

import "strings"

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
