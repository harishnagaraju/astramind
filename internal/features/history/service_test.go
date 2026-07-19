package history

import (
	"testing"

	"github.com/harishnagaraju/astramind/internal/infrastructure/models"
	"github.com/harishnagaraju/astramind/internal/infrastructure/storage"
)

// newTestService returns a Service backed by a temp-directory store,
// isolated from the real data/sessions folder and automatically
// cleaned up by the test framework, even if the test fails partway
// through.
func newTestService(t *testing.T) *Service {
	t.Helper()

	store := storage.NewFileHistoryStore(t.TempDir())

	return NewServiceWithStore(store)
}

func TestService_SaveAndLoad(t *testing.T) {
	service := newTestService(t)

	messages := []models.Message{
		{Role: "user", Content: "hello"},
		{Role: "assistant", Content: "hi there"},
	}

	if err := service.Save("test-session", messages); err != nil {
		t.Fatal(err)
	}

	loaded, err := service.Load("test-session")
	if err != nil {
		t.Fatal(err)
	}

	if len(loaded) != len(messages) {
		t.Fatalf("expected %d messages, got %d", len(messages), len(loaded))
	}

	if loaded[0].Content != "hello" {
		t.Fatalf("unexpected content: %s", loaded[0].Content)
	}
}

func TestService_LoadMissingSession(t *testing.T) {
	service := newTestService(t)

	loaded, err := service.Load("does-not-exist")
	if err != nil {
		t.Fatal(err)
	}

	if len(loaded) != 0 {
		t.Fatalf("expected empty history, got %d messages", len(loaded))
	}
}

func TestService_ListSessions(t *testing.T) {
	service := newTestService(t)

	messages := []models.Message{{Role: "user", Content: "hi"}}

	if err := service.Save("session-a", messages); err != nil {
		t.Fatal(err)
	}

	if err := service.Save("session-b", messages); err != nil {
		t.Fatal(err)
	}

	sessions, err := service.ListSessions()
	if err != nil {
		t.Fatal(err)
	}

	if len(sessions) != 2 {
		t.Fatalf("expected 2 sessions, got %d", len(sessions))
	}
}

func TestService_Delete(t *testing.T) {
	service := newTestService(t)

	messages := []models.Message{{Role: "user", Content: "hi"}}

	if err := service.Save("to-delete", messages); err != nil {
		t.Fatal(err)
	}

	if err := service.Delete("to-delete"); err != nil {
		t.Fatal(err)
	}

	sessions, err := service.ListSessions()
	if err != nil {
		t.Fatal(err)
	}

	for _, s := range sessions {
		if s == "to-delete" {
			t.Fatal("expected session to be deleted")
		}
	}
}

func TestNewService_UsesRealDataDirectory(t *testing.T) {
	// NewService() (no store injected) must keep working exactly as
	// before this change - every existing caller (bootstrap.go,
	// session.Service, search.Service) depends on this behavior.
	service := NewService()

	if service.store == nil {
		t.Fatal("expected NewService() to have a default store configured")
	}
}
