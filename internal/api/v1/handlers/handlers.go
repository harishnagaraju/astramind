package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/harishnagaraju/astramind/internal/api/v1/errors"
	"github.com/harishnagaraju/astramind/internal/api/v1/request"
	"github.com/harishnagaraju/astramind/internal/api/v1/response"
	"github.com/harishnagaraju/astramind/internal/features/kb"
	"github.com/harishnagaraju/astramind/internal/infrastructure/ai"
	"github.com/harishnagaraju/astramind/internal/infrastructure/config"
	"github.com/harishnagaraju/astramind/internal/infrastructure/models"
)

type Config struct {
	ProviderName    string
	Model           string
	Version         string
	APIKey          string
	ProviderManager *ai.ProviderManager
	KnowledgeBase   *kb.Manager
}

type Handler struct {
	config Config
}

func New(config Config) *Handler { return &Handler{config: config} }

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.methodNotAllowed(w, r)
		return
	}
	writeJSON(w, http.StatusOK, response.Health{Status: "ok"})
}

func (h *Handler) Status(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.methodNotAllowed(w, r)
		return
	}
	provider := h.config.ProviderName
	if provider == "" && h.config.ProviderManager != nil {
		provider = h.config.ProviderManager.ProviderName()
	}
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
		h.methodNotAllowed(w, r)
		return
	}
	writeJSON(w, http.StatusOK, response.Version{Version: h.config.Version})
}

type chatRequest struct {
	Model       string           `json:"model,omitempty"`
	Messages    []models.Message `json:"messages"`
	Temperature *float64         `json:"temperature,omitempty"`
}

type chatResponse struct {
	Message models.Message `json:"message"`
	Model   string         `json:"model,omitempty"`
}

func (h *Handler) Chat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.methodNotAllowed(w, r)
		return
	}
	if h.config.ProviderManager == nil {
		h.serverError(w, r, "provider manager is not configured")
		return
	}

	var req chatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.badRequest(w, r, "invalid request body")
		return
	}
	if len(req.Messages) == 0 {
		h.badRequest(w, r, "messages is required")
		return
	}
	if len(req.Messages) > config.MaxMessages {
		h.badRequest(w, r, fmt.Sprintf("maximum %d messages allowed", config.MaxMessages))
		return
	}

	model := req.Model
	if model == "" {
		model = h.config.Model
	}

	reply, err := h.config.ProviderManager.Chat(ai.ChatRequest{
		Model:       model,
		APIKey:      h.config.APIKey,
		Messages:    req.Messages,
		Temperature: req.Temperature,
	})
	if err != nil {
		h.serverError(w, r, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, chatResponse{
		Message: models.Message{Role: "assistant", Content: reply},
		Model:   model,
	})
}

func (h *Handler) ChatStream(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.methodNotAllowed(w, r)
		return
	}
	if h.config.ProviderManager == nil {
		h.serverError(w, r, "provider manager is not configured")
		return
	}

	var req chatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.badRequest(w, r, "invalid request body")
		return
	}
	if len(req.Messages) == 0 {
		h.badRequest(w, r, "messages is required")
		return
	}
	if len(req.Messages) > config.MaxMessages {
		h.badRequest(w, r, fmt.Sprintf("maximum %d messages allowed", config.MaxMessages))
		return
	}

	model := req.Model
	if model == "" {
		model = h.config.Model
	}

	stream, err := h.config.ProviderManager.Stream(r.Context(), ai.ChatRequest{
		Model:       model,
		APIKey:      h.config.APIKey,
		Messages:    req.Messages,
		Temperature: req.Temperature,
	})
	if err != nil {
		h.serverError(w, r, err.Error())
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	flusher, ok := w.(http.Flusher)
	if !ok {
		h.serverError(w, r, "streaming is not supported by the HTTP server")
		return
	}

	for event := range stream.Events() {
		switch event.Type {
		case ai.StreamEventToken:
			writeSSE(w, "token", map[string]string{"content": event.Content})
		case ai.StreamEventDone:
			writeSSE(w, "done", map[string]string{})
		case ai.StreamEventError:
			message := "stream error"
			if event.Err != nil {
				message = event.Err.Error()
			}
			writeSSE(w, "error", map[string]string{"message": message})
		}
		flusher.Flush()
	}
}

type documentResponse struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	ChunkCount int    `json:"chunk_count"`
}

func (h *Handler) Documents(w http.ResponseWriter, r *http.Request) {
	if h.config.KnowledgeBase == nil {
		h.serverError(w, r, "knowledge base is not configured")
		return
	}
	switch r.Method {
	case http.MethodGet:
		documents, err := h.config.KnowledgeBase.ListKnowledge()
		if err != nil {
			h.serverError(w, r, err.Error())
			return
		}
		result := make([]documentResponse, 0, len(documents))
		for _, doc := range documents {
			result = append(result, documentResponse{ID: doc.ID, Name: doc.Name, ChunkCount: doc.ChunkCount})
		}
		writeJSON(w, http.StatusOK, result)
	case http.MethodPost:
		h.importDocument(w, r)
	default:
		h.methodNotAllowed(w, r)
	}
}

func (h *Handler) Document(w http.ResponseWriter, r *http.Request) {
	if h.config.KnowledgeBase == nil {
		h.serverError(w, r, "knowledge base is not configured")
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/documents/")
	if id == "" || strings.Contains(id, "/") {
		h.badRequest(w, r, "document id is required")
		return
	}

	switch r.Method {
	case http.MethodGet:
		doc, err := h.config.KnowledgeBase.GetKnowledge(id)
		if err != nil {
			h.notFound(w, r, "document not found")
			return
		}
		writeJSON(w, http.StatusOK, doc)
	case http.MethodDelete:
		if err := h.config.KnowledgeBase.RemoveKnowledge(id); err != nil {
			h.notFound(w, r, "document not found")
			return
		}
		w.WriteHeader(http.StatusNoContent)
	default:
		h.methodNotAllowed(w, r)
	}
}

func (h *Handler) importDocument(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	if err != nil {
		h.badRequest(w, r, "file is required")
		return
	}
	defer file.Close()

	doc, err := saveAndImport(h.config.KnowledgeBase, file, header)
	if err != nil {
		h.badRequest(w, r, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, documentResponse{
		ID: doc.ID, Name: doc.Name, ChunkCount: doc.ChunkCount,
	})
}

func (h *Handler) AnalyzeDocument(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.methodNotAllowed(w, r)
		return
	}
	if h.config.KnowledgeBase == nil || h.config.ProviderManager == nil {
		h.serverError(w, r, "document analysis is not configured")
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/documents/")
	id = strings.TrimSuffix(id, "/analyze")
	if id == "" {
		h.badRequest(w, r, "document id is required")
		return
	}

	doc, err := h.config.KnowledgeBase.GetKnowledge(id)
	if err != nil {
		h.notFound(w, r, "document not found")
		return
	}

	var body struct {
		Prompt string `json:"prompt,omitempty"`
	}
	if r.Body != nil && r.ContentLength != 0 {
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			h.badRequest(w, r, "invalid request body")
			return
		}
	}
	prompt := body.Prompt
	if prompt == "" {
		prompt = "Analyze the following document and provide a concise summary of its key points.\n\n" + doc.Content
	} else {
		prompt += "\n\nDocument:\n" + doc.Content
	}

	reply, err := h.config.ProviderManager.Chat(ai.ChatRequest{
		Model:  h.config.Model,
		APIKey: h.config.APIKey,
		Messages: []models.Message{{Role: "user", Content: prompt}},
	})
	if err != nil {
		h.serverError(w, r, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"document_id": id, "analysis": reply})
}

type embeddingRequest struct {
	Text  string `json:"text"`
	Model string `json:"model,omitempty"`
}

func (h *Handler) Embeddings(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.methodNotAllowed(w, r)
		return
	}
	if h.config.ProviderManager == nil {
		h.serverError(w, r, "provider manager is not configured")
		return
	}
	var req embeddingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.badRequest(w, r, "invalid request body")
		return
	}
	if strings.TrimSpace(req.Text) == "" {
		h.badRequest(w, r, "text is required")
		return
	}
	embedding, err := h.config.ProviderManager.Embed(ai.EmbeddingRequest{
		Model: req.Model, APIKey: h.config.APIKey, Text: req.Text,
	})
	if err != nil {
		h.serverError(w, r, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"model": req.Model, "embedding": embedding})
}

type searchRequest struct {
	Query string `json:"query"`
	Mode  string `json:"mode,omitempty"`
}

func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.methodNotAllowed(w, r)
		return
	}
	if h.config.KnowledgeBase == nil {
		h.serverError(w, r, "knowledge base is not configured")
		return
	}
	var req searchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.badRequest(w, r, "invalid request body")
		return
	}
	if strings.TrimSpace(req.Query) == "" {
		h.badRequest(w, r, "query is required")
		return
	}

	if strings.EqualFold(req.Mode, "semantic") {
		results, err := h.config.KnowledgeBase.SemanticSearch(req.Query)
		if err != nil {
			h.serverError(w, r, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"mode": "semantic", "results": results})
		return
	}

	results, err := h.config.KnowledgeBase.Search(req.Query)
	if err != nil {
		h.serverError(w, r, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"mode": "keyword", "results": results})
}

func (h *Handler) Providers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.methodNotAllowed(w, r)
		return
	}
	provider := h.config.ProviderName
	if provider == "" && h.config.ProviderManager != nil {
		provider = h.config.ProviderManager.ProviderName()
	}
	writeJSON(w, http.StatusOK, map[string]any{"providers": []string{provider}})
}

func (h *Handler) Models(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.methodNotAllowed(w, r)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"models": []string{h.config.Model}})
}

func saveAndImport(manager *kb.Manager, file multipart.File, header *multipart.FileHeader) (*kb.Document, error) {
	name := filepath.Base(header.Filename)
	if name == "." || name == "" {
		return nil, fmt.Errorf("invalid filename")
	}
	ext := strings.ToLower(filepath.Ext(name))
	if ext != ".txt" && ext != ".md" && ext != ".docx" {
		return nil, fmt.Errorf("unsupported document type")
	}
	if err := os.MkdirAll(filepath.Join("data", "uploads"), 0755); err != nil {
		return nil, err
	}
	path := filepath.Join("data", "uploads", name)
	dest, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(dest, file); err != nil {
		dest.Close()
		return nil, err
	}
	if err := dest.Close(); err != nil {
		return nil, err
	}
	return manager.ImportDocument(path)
}

func writeSSE(w http.ResponseWriter, event string, data any) {
	payload, _ := json.Marshal(data)
	fmt.Fprintf(w, "event: %s\ndata: %s\n\n", event, payload)
}

func (h *Handler) methodNotAllowed(w http.ResponseWriter, r *http.Request) {
	errors.Write(w, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed", request.ID(r.Context()))
}
func (h *Handler) badRequest(w http.ResponseWriter, r *http.Request, message string) {
	errors.Write(w, http.StatusBadRequest, "bad_request", message, request.ID(r.Context()))
}
func (h *Handler) notFound(w http.ResponseWriter, r *http.Request, message string) {
	errors.Write(w, http.StatusNotFound, "not_found", message, request.ID(r.Context()))
}
func (h *Handler) serverError(w http.ResponseWriter, r *http.Request, message string) {
	errors.Write(w, http.StatusInternalServerError, "internal_error", message, request.ID(r.Context()))
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

var _ = strconv.Itoa
