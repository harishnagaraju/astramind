package engine

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/harishnagaraju/astramind/internal/features/kb"
	"github.com/harishnagaraju/astramind/internal/infrastructure/ai"
)

// webTestEmbedder is a minimal, deterministic test double for
// kb.Embedder, used so these tests never touch a real embedding
// provider.
type webTestEmbedder struct{}

func (webTestEmbedder) Embed(text string) ([]float32, error) {
	return []float32{0.1, 0.2, 0.3}, nil
}

// newTestApp returns an App with KnowledgeBase and ProviderManager
// wired directly (bypassing initialize(), which needs a real .env),
// backed by an isolated temp directory so these tests never touch
// the real data/ folder.
func newTestApp(t *testing.T) *App {
	t.Helper()

	tempDir := t.TempDir()

	storage := kb.NewJSONStorage(tempDir)
	manager := kb.NewManager(storage)
	manager.SetEmbedder(webTestEmbedder{})

	app := &App{
		providerName: "mock",
	}

	app.deps.KnowledgeBase = manager
	app.deps.ProviderManager = ai.NewProviderManager(&ai.MockProvider{})

	return app
}

func TestHandleAPIStatus(t *testing.T) {
	app := newTestApp(t)

	req := httptest.NewRequest(http.MethodGet, "/api/status", nil)
	rec := httptest.NewRecorder()

	app.handleAPIStatus(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var resp statusResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}

	if resp.Provider != "mock" {
		t.Fatalf("expected provider 'mock', got %q", resp.Provider)
	}
}

func TestHandleAPIDocuments_EmptyInitially(t *testing.T) {
	app := newTestApp(t)

	req := httptest.NewRequest(http.MethodGet, "/api/documents", nil)
	rec := httptest.NewRecorder()

	app.handleAPIDocuments(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var docs []documentSummary
	if err := json.NewDecoder(rec.Body).Decode(&docs); err != nil {
		t.Fatal(err)
	}

	if len(docs) != 0 {
		t.Fatalf("expected 0 documents, got %d", len(docs))
	}
}

func TestHandleAPIDocuments_ImportAndList(t *testing.T) {
	app := newTestApp(t)

	// Build a multipart/form-data upload request, the same shape the
	// browser's FormData would send.
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	part, err := writer.CreateFormFile("file", "sample.txt")
	if err != nil {
		t.Fatal(err)
	}
	part.Write([]byte("hello world, this is a test document"))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/documents", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rec := httptest.NewRecorder()

	app.handleAPIDocuments(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	// Now list and confirm it shows up.
	listReq := httptest.NewRequest(http.MethodGet, "/api/documents", nil)
	listRec := httptest.NewRecorder()

	app.handleAPIDocuments(listRec, listReq)

	var docs []documentSummary
	if err := json.NewDecoder(listRec.Body).Decode(&docs); err != nil {
		t.Fatal(err)
	}

	if len(docs) != 1 {
		t.Fatalf("expected 1 document, got %d", len(docs))
	}

	if docs[0].Name != "sample.txt" {
		t.Fatalf("expected sample.txt, got %s", docs[0].Name)
	}

	// Clean up the upload artifact this test created outside t.TempDir().
	os.Remove(filepath.Join("data", "uploads", "sample.txt"))
}

func TestHandleAPIAsk_NoDocuments(t *testing.T) {
	app := newTestApp(t)

	body, _ := json.Marshal(askRequest{Question: "anything"})
	req := httptest.NewRequest(http.MethodPost, "/api/ask", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	app.handleAPIAsk(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var resp askResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}

	if resp.Answer == "" {
		t.Fatal("expected a fallback answer explaining no knowledge was found")
	}
}

func TestHandleAPIAsk_EmptyQuestion(t *testing.T) {
	app := newTestApp(t)

	body, _ := json.Marshal(askRequest{Question: ""})
	req := httptest.NewRequest(http.MethodPost, "/api/ask", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	app.handleAPIAsk(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestHandleAPIAsk_WithDocuments(t *testing.T) {
	app := newTestApp(t)

	tempDir := t.TempDir()
	source := filepath.Join(tempDir, "sample.txt")

	if err := os.WriteFile(source, []byte("hello world"), 0644); err != nil {
		t.Fatal(err)
	}

	if _, err := app.deps.KnowledgeBase.ImportDocument(source); err != nil {
		t.Fatal(err)
	}

	reqBody, _ := json.Marshal(askRequest{Question: "what does the document say"})
	req := httptest.NewRequest(http.MethodPost, "/api/ask", bytes.NewReader(reqBody))
	rec := httptest.NewRecorder()

	app.handleAPIAsk(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var resp askResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}

	if len(resp.Sources) == 0 {
		t.Fatal("expected at least one cited source")
	}
}
