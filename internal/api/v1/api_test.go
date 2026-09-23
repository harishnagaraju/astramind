package v1_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	v1 "github.com/harishnagaraju/astramind/internal/api/v1"
)

func TestPlatformAPIHealth(t *testing.T) {
	handler := v1.New(v1.Config{
		ProviderName: "mock",
		Model:        "test-model",
		Version:      "v0.9.2",
	})

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if got := rec.Header().Get("X-Request-ID"); got == "" {
		t.Fatal("expected X-Request-ID response header")
	}
	if got := rec.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("expected application/json content type, got %q", got)
	}
}

func TestPlatformAPIStatus(t *testing.T) {
	handler := v1.New(v1.Config{
		ProviderName: "ollama",
		Model:        "gemma3",
		Version:      "v0.9.2",
	})

	req := httptest.NewRequest(http.MethodGet, "/status", nil)
	req.Header.Set("X-Request-ID", "test-request-id")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if got := rec.Header().Get("X-Request-ID"); got != "test-request-id" {
		t.Fatalf("expected propagated request id, got %q", got)
	}
}

func TestPlatformAPIVersion(t *testing.T) {
	handler := v1.New(v1.Config{Version: "v1.2.3"})

	req := httptest.NewRequest(http.MethodGet, "/version", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestPlatformAPIMethodNotAllowed(t *testing.T) {
	handler := v1.New(v1.Config{Version: "v0.9.2"})

	req := httptest.NewRequest(http.MethodPost, "/health", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rec.Code)
	}
}
