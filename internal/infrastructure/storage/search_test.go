package storage

import (
	"testing"

	"github.com/harishnagaraju/astramind/internal/infrastructure/models"
)

func TestSearchMessages_EmptyQuery(t *testing.T) {
	messages := []models.Message{
		{
			Role:    "user",
			Content: "Hello World",
		},
	}

	results := SearchMessages(messages, "")

	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

func TestSearchMessages_NoMatch(t *testing.T) {
	messages := []models.Message{
		{
			Role:    "user",
			Content: "Hello World",
		},
	}

	results := SearchMessages(messages, "python")

	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

func TestSearchMessages_SingleMatch(t *testing.T) {
	messages := []models.Message{
		{
			Role:    "user",
			Content: "Hello World",
		},
		{
			Role:    "assistant",
			Content: "Go is an excellent language.",
		},
	}

	results := SearchMessages(messages, "Go")

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}

	if results[0].Index != 1 {
		t.Fatalf("expected index 1, got %d", results[0].Index)
	}

	if results[0].Role != "assistant" {
		t.Fatalf("expected assistant, got %s", results[0].Role)
	}

	if results[0].Content != "Go is an excellent language." {
		t.Fatalf("unexpected content")
	}
}

func TestSearchMessages_MultipleMatches(t *testing.T) {
	messages := []models.Message{
		{
			Role:    "user",
			Content: "Go is fast.",
		},
		{
			Role:    "assistant",
			Content: "Python is popular.",
		},
		{
			Role:    "assistant",
			Content: "Go supports concurrency.",
		},
	}

	results := SearchMessages(messages, "go")

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestSearchMessages_CaseInsensitive(t *testing.T) {
	messages := []models.Message{
		{
			Role:    "assistant",
			Content: "GOLANG",
		},
	}

	results := SearchMessages(messages, "golang")

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
}

func TestSearchMessages_UserMessage(t *testing.T) {
	messages := []models.Message{
		{
			Role:    "user",
			Content: "Learning Go",
		},
	}

	results := SearchMessages(messages, "Learning")

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}

	if results[0].Role != "user" {
		t.Fatalf("expected user role")
	}
}

func TestSearchMessages_AssistantMessage(t *testing.T) {
	messages := []models.Message{
		{
			Role:    "assistant",
			Content: "Go Modules",
		},
	}

	results := SearchMessages(messages, "Modules")

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}

	if results[0].Role != "assistant" {
		t.Fatalf("expected assistant role")
	}
}

func TestSearchMessages_SubstringMatch(t *testing.T) {
	messages := []models.Message{
		{
			Role:    "assistant",
			Content: "Programming Language",
		},
	}

	results := SearchMessages(messages, "gram")

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
}
