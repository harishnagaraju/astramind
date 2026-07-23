package kb

import (
	"strings"
	"testing"
)

func TestBuildPromptSingleResult(t *testing.T) {
	results := []SearchResult{
		{
			DocumentID: "go-guide",
			Content:    "Go is a programming language.",
		},
	}

	prompt := BuildPrompt(
		"What is Go?",
		results,
	)

	if !strings.Contains(prompt, "Go is a programming language.") {
		t.Fatal("expected document content in prompt")
	}

	if !strings.Contains(prompt, "What is Go?") {
		t.Fatal("expected question in prompt")
	}
}

func TestBuildPromptMultipleResults(t *testing.T) {
	results := []SearchResult{
		{
			DocumentID: "doc1",
			Content:    "First document.",
		},
		{
			DocumentID: "doc2",
			Content:    "Second document.",
		},
	}

	prompt := BuildPrompt(
		"Example question",
		results,
	)

	if !strings.Contains(prompt, "First document.") {
		t.Fatal("missing first document")
	}

	if !strings.Contains(prompt, "Second document.") {
		t.Fatal("missing second document")
	}
}

func TestBuildPromptEmptyResults(t *testing.T) {
	prompt := BuildPrompt(
		"What is Go?",
		nil,
	)

	if !strings.Contains(prompt, "Question:") {
		t.Fatal("expected question section")
	}
}

func TestBuildPromptEmptyQuestion(t *testing.T) {
	results := []SearchResult{
		{
			DocumentID: "doc",
			Content:    "Knowledge",
		},
	}

	prompt := BuildPrompt(
		"",
		results,
	)

	if !strings.Contains(prompt, "Knowledge") {
		t.Fatal("expected knowledge content")
	}
}

// TestBuildSemanticPromptForcesPerSourceEnumeration guards against a
// regression to the previous single-trailing-instruction approach.
// v0.9.1 validation testing showed a trailing "include everything"
// instruction alone was not sufficient to prevent the model from
// narrowly matching only the question's literal keywords - the
// prompt must explicitly number the sources and instruct the model
// to work through each one individually.
func TestBuildSemanticPromptForcesPerSourceEnumeration(t *testing.T) {
	results := []SemanticSearchResult{
		{DocumentID: "doc1", Content: "First entry."},
		{DocumentID: "doc2", Content: "Second entry."},
		{DocumentID: "doc3", Content: "Third entry."},
	}

	prompt := BuildSemanticPrompt("What are all the entries?", results)

	if !strings.Contains(prompt, "There are 3 sources above, numbered 1 through 3") {
		t.Fatal("expected prompt to explicitly state the source count for per-source enumeration")
	}

	if !strings.Contains(prompt, "Work through each source in order") {
		t.Fatal("expected prompt to instruct working through sources individually")
	}

	if !strings.Contains(prompt, "even if that source doesn't use the same words as the question") {
		t.Fatal("expected prompt to explicitly address the narrow-keyword-matching failure mode")
	}
}
