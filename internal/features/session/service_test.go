package session

import (
	"testing"

	"github.com/harishnagaraju/astramind/internal/features/history"
	"github.com/harishnagaraju/astramind/internal/infrastructure/models"
	"github.com/harishnagaraju/astramind/internal/infrastructure/storage"
)

// newTestService returns a Service backed by a temp-directory store,
// isolated from the real data/sessions folder and automatically
// cleaned up, even if the test fails partway through.
func newTestService(t *testing.T) *Service {
	t.Helper()

	fileStore := storage.NewFileHistoryStore(t.TempDir())
	historySvc := history.NewServiceWithStore(fileStore)

	return NewServiceWithStore(historySvc, fileStore)
}

func TestService_CreateAndExists(t *testing.T) {
	service := newTestService(t)

	if service.Exists("new-session") {
		t.Fatal("expected session not to exist yet")
	}

	if err := service.Create("new-session"); err != nil {
		t.Fatal(err)
	}

	if !service.Exists("new-session") {
		t.Fatal("expected session to exist after Create")
	}
}

func TestService_CreateIsIdempotent(t *testing.T) {
	service := newTestService(t)

	if err := service.Create("repeat-session"); err != nil {
		t.Fatal(err)
	}

	// Creating an already-existing session should not error or wipe
	// its content.
	messages := []models.Message{{Role: "user", Content: "hello"}}

	if err := service.Save("repeat-session", messages); err != nil {
		t.Fatal(err)
	}

	if err := service.Create("repeat-session"); err != nil {
		t.Fatal(err)
	}

	loaded, err := service.Load("repeat-session")
	if err != nil {
		t.Fatal(err)
	}

	if len(loaded) != 1 {
		t.Fatalf("expected Create on an existing session to leave content untouched, got %d messages", len(loaded))
	}
}

func TestService_SaveLoadDeleteList(t *testing.T) {
	service := newTestService(t)

	messages := []models.Message{{Role: "user", Content: "hi"}}

	if err := service.Save("session-x", messages); err != nil {
		t.Fatal(err)
	}

	loaded, err := service.Load("session-x")
	if err != nil {
		t.Fatal(err)
	}

	if len(loaded) != 1 {
		t.Fatalf("expected 1 message, got %d", len(loaded))
	}

	sessions, err := service.List()
	if err != nil {
		t.Fatal(err)
	}

	found := false
	for _, s := range sessions {
		if s == "session-x" {
			found = true
		}
	}
	if !found {
		t.Fatal("expected session-x to appear in List()")
	}

	if err := service.Delete("session-x"); err != nil {
		t.Fatal(err)
	}

	if service.Exists("session-x") {
		t.Fatal("expected session-x to no longer exist after Delete")
	}
}
