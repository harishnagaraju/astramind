package chat

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/harishnagaraju/astramind/internal/features/kb"
	"github.com/harishnagaraju/astramind/internal/infrastructure/ai"
)

// fakeEmbedder is a minimal test double satisfying kb.Embedder.
type fakeEmbedder struct{}

func (f *fakeEmbedder) Embed(text string) ([]float32, error) {
	return []float32{0.1, 0.2, 0.3}, nil
}

func TestHandleKBSemanticSearchInvalidArguments(t *testing.T) {
	service := newTestService(t)

	err := service.handleKBSemanticSearch([]string{
		"/kb",
		"ssearch",
	})

	if err == nil {
		t.Fatal("expected error for missing query")
	}
}

func TestHandleKBSemanticSearchNoEmbedderConfigured(t *testing.T) {
	service := newTestService(t)

	// newTestService does not configure an embedder - should
	// fail loudly rather than silently returning no results.
	err := service.handleKBSemanticSearch([]string{
		"/kb",
		"ssearch",
		"hello",
	})

	if err == nil {
		t.Fatal("expected error when no embedder is configured")
	}
}

func TestHandleKBSemanticSearchSuccess(t *testing.T) {

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

	handled, err := service.HandleKnowledgeCommand("/kb ssearch hello")
	if err != nil {
		t.Fatal(err)
	}

	if !handled {
		t.Fatal("expected command to be handled")
	}
}
