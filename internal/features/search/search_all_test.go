package search

import (
	"testing"

	"github.com/harishnagaraju/astramind/internal/features/history"
	"github.com/harishnagaraju/astramind/internal/infrastructure/models"
)

func TestSearchAllSessions(t *testing.T) {

	historySvc := history.NewService()

	session1 := "search-all-session-1"
	session2 := "search-all-session-2"

	err := historySvc.Save(session1, []models.Message{
		{
			Role:    "user",
			Content: "Learning Go",
		},
		{
			Role:    "assistant",
			Content: "Go is fast.",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	err = historySvc.Save(session2, []models.Message{
		{
			Role:    "user",
			Content: "Python Programming",
		},
		{
			Role:    "assistant",
			Content: "Go Modules",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	results, err := SearchAllSessions(historySvc, "go")
	if err != nil {
		t.Fatal(err)
	}

	if len(results) < 3 {
		t.Fatalf("expected at least 3 matches, got %d", len(results))
	}

	foundSession1 := false
	foundSession2 := false

	for _, result := range results {

		switch result.Session {

		case session1:
			foundSession1 = true

		case session2:
			foundSession2 = true
		}
	}

	if !foundSession1 {
		t.Fatal("expected results from session1")
	}

	if !foundSession2 {
		t.Fatal("expected results from session2")
	}

	_ = historySvc.Delete(session1)
	_ = historySvc.Delete(session2)
}
