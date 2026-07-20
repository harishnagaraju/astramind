package chat

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/harishnagaraju/astramind/internal/features/kb"
	"github.com/harishnagaraju/astramind/internal/infrastructure/ai"
)

func TestHandleKBAskInvalidArguments(t *testing.T) {
	service := newTestService(t)

	err := service.handleKBAsk([]string{
		"/kb",
		"ask",
	})

	if err == nil {
		t.Fatal("expected error for missing question")
	}
}

func TestHandleKBAskNoRelevantKnowledge(t *testing.T) {

	providerManager := ai.NewProviderManager(&ai.MockProvider{})

	tempDir := t.TempDir()

	storage := kb.NewJSONStorage(tempDir)
	kbManager := kb.NewManager(storage)
	kbManager.SetEmbedder(&fakeEmbedder{})

	service := NewService(Dependencies{
		ProviderManager: providerManager,
		KnowledgeBase:   kbManager,
	})

	// No documents imported - SemanticSearch should return nothing.
	handled, err := service.HandleKnowledgeCommand("/kb ask anything")
	if err != nil {
		t.Fatal(err)
	}

	if !handled {
		t.Fatal("expected command to be handled")
	}
}

func TestHandleKBAskSuccess(t *testing.T) {

	providerManager := ai.NewProviderManager(&ai.MockProvider{})

	tempDir := t.TempDir()

	storage := kb.NewJSONStorage(tempDir)
	kbManager := kb.NewManager(storage)
	kbManager.SetEmbedder(&fakeEmbedder{})

	service := NewService(Dependencies{
		ProviderManager: providerManager,
		KnowledgeBase:   kbManager,
	})

	source := filepath.Join(tempDir, "sample.txt")

	if err := os.WriteFile(source, []byte("hello world"), 0644); err != nil {
		t.Fatal(err)
	}

	if _, err := kbManager.ImportDocument(source); err != nil {
		t.Fatal(err)
	}

	handled, err := service.HandleKnowledgeCommand("/kb ask what does the document say")
	if err != nil {
		t.Fatal(err)
	}

	if !handled {
		t.Fatal("expected command to be handled")
	}
}
