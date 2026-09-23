package v1_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	v1 "github.com/harishnagaraju/astramind/internal/api/v1"
)

func TestMount(t *testing.T) {
	mux := http.NewServeMux()
	v1.Mount(mux, v1.Config{ProviderName: "mock", Version: "v0.9.2"})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}
