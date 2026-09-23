package v1_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	platformapi "github.com/harishnagaraju/astramind/internal/api/v1"
	"github.com/harishnagaraju/astramind/internal/infrastructure/ai"
)

func testHandler() http.Handler {
	manager := ai.NewProviderManager(&ai.MockProvider{})
	return platformapi.New(platformapi.Config{
		ProviderName:    "mock",
		Model:           "mock-model",
		Version:         "v1.0.0",
		APIKey:          "",
		ProviderManager: manager,
	})
}

func TestHealthAndRequestID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	req.Header.Set("X-Request-ID", "test-request")
	rec := httptest.NewRecorder()

	testHandler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if got := rec.Header().Get("X-Request-ID"); got != "test-request" {
		t.Fatalf("expected request id propagation, got %q", got)
	}
	if !strings.Contains(rec.Body.String(), `"status":"ok"`) {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}
}

func TestChat(t *testing.T) {
	body := strings.NewReader(`{"messages":[{"role":"user","content":"hello"}]}`)
	req := httptest.NewRequest(http.MethodPost, "/chat", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	testHandler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var result struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}
	if result.Message.Role != "assistant" || result.Message.Content == "" {
		t.Fatalf("unexpected response: %+v", result)
	}
}

func TestChatStream(t *testing.T) {
	body := strings.NewReader(`{"messages":[{"role":"user","content":"hello"}]}`)
	req := httptest.NewRequest(http.MethodPost, "/chat/stream", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	testHandler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "event: token") {
		t.Fatalf("missing token event: %s", rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "event: done") {
		t.Fatalf("missing done event: %s", rec.Body.String())
	}
}

func TestEmbeddings(t *testing.T) {
	body := strings.NewReader(`{"text":"hello world"}`)
	req := httptest.NewRequest(http.MethodPost, "/embeddings", body)
	rec := httptest.NewRecorder()

	testHandler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "embedding") {
		t.Fatalf("missing embedding response: %s", rec.Body.String())
	}
}
