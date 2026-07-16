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
