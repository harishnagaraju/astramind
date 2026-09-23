package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/harishnagaraju/astramind/internal/api/v1/errors"
	"github.com/harishnagaraju/astramind/internal/api/v1/request"
	"github.com/harishnagaraju/astramind/internal/api/v1/response"
)

type Config struct {
	ProviderName string
	Model        string
	Version      string
}

type Handler struct {
	config Config
}

func New(config Config) *Handler {
	return &Handler{config: config}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errors.Write(w, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed", request.ID(r.Context()))
		return
	}
	writeJSON(w, http.StatusOK, response.Health{Status: "ok"})
}

func (h *Handler) Status(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errors.Write(w, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed", request.ID(r.Context()))
		return
	}

	provider := h.config.ProviderName
	if provider == "" {
		provider = "unknown"
	}

	writeJSON(w, http.StatusOK, response.Status{
		Status:   "ok",
		Provider: provider,
		Model:    h.config.Model,
	})
}

func (h *Handler) Version(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errors.Write(w, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed", request.ID(r.Context()))
		return
	}
	writeJSON(w, http.StatusOK, response.Version{Version: h.config.Version})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
