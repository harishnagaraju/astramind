package engine

import (
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/harishnagaraju/astramind/internal/features/kb"
	"github.com/harishnagaraju/astramind/internal/infrastructure/ai"
	"github.com/harishnagaraju/astramind/internal/infrastructure/models"
)

//go:embed webui/index.html
var webUIFiles embed.FS

// askRequest is the JSON body for POST /api/ask.
type askRequest struct {
	Question string `json:"question"`
}

// askResponse is the JSON body returned by POST /api/ask.
type askResponse struct {
	Answer  string   `json:"answer"`
	Sources []string `json:"sources"`
	Error   string   `json:"error,omitempty"`
}

// documentSummary is the JSON shape returned by GET /api/documents.
type documentSummary struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	ChunkCount int    `json:"chunk_count"`
}

// statusResponse is the JSON body returned by GET /api/status.
type statusResponse struct {
	Provider string `json:"provider"`
}

// runWeb starts the local web server and serves the embedded UI on
// the given address (e.g. "localhost:8420"). It reuses a.deps
// directly - the same services already wired for interactive/script
// mode - so behavior (including which AI provider is active) is
// identical to the CLI.
func (a *App) runWeb(addr string) error {

	mux := http.NewServeMux()

	indexHTML, err := webUIFiles.ReadFile("webui/index.html")
	if err != nil {
		return err
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(indexHTML)
	})

	mux.HandleFunc("/api/status", a.handleAPIStatus)
	mux.HandleFunc("/api/documents", a.handleAPIDocuments)
	mux.HandleFunc("/api/ask", a.handleAPIAsk)

	fmt.Printf("AstraMind web UI running at http://%s\n", addr)
	fmt.Println("Press Ctrl+C to stop.")

	return http.ListenAndServe(addr, mux)
}

func (a *App) handleAPIStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(statusResponse{
		Provider: a.providerName,
	})
}

func (a *App) handleAPIDocuments(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case http.MethodGet:
		a.listDocuments(w, r)

	case http.MethodPost:
		a.importDocument(w, r)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (a *App) listDocuments(w http.ResponseWriter, r *http.Request) {

	documents, err := a.deps.KnowledgeBase.ListKnowledge()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	summaries := make([]documentSummary, 0, len(documents))

	for _, doc := range documents {
		summaries = append(summaries, documentSummary{
			ID:         doc.ID,
			Name:       doc.Name,
			ChunkCount: doc.ChunkCount,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summaries)
}

func (a *App) importDocument(w http.ResponseWriter, r *http.Request) {

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "no file provided", http.StatusBadRequest)
		return
	}
	defer file.Close()

	uploadDir := filepath.Join("data", "uploads")

	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	destPath := filepath.Join(uploadDir, header.Filename)

	dest, err := os.Create(destPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dest.Close()

	if _, err := io.Copy(dest, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Reuses the exact same ImportDocument used by the CLI's /kb
	// import - same chunking, same embedding behavior, no duplicated
	// logic.
	if _, err := a.deps.KnowledgeBase.ImportDocument(destPath); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *App) handleAPIAsk(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req askRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if req.Question == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(askResponse{Error: "question is required"})
		return
	}

	results, err := a.deps.KnowledgeBase.SemanticSearch(req.Question)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(askResponse{Error: err.Error()})
		return
	}

	if len(results) == 0 {
		json.NewEncoder(w).Encode(askResponse{
			Answer:  "No relevant knowledge found to answer this question.",
			Sources: []string{},
		})
		return
	}

	prompt := kb.BuildSemanticPrompt(req.Question, results)

	reply, err := a.deps.ProviderManager.Chat(ai.ChatRequest{
		Messages: []models.Message{
			{Role: "user", Content: prompt},
		},
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(askResponse{Error: err.Error()})
		return
	}

	seen := make(map[string]bool)
	sources := make([]string, 0)

	for _, result := range results {
		if seen[result.DocumentID] {
			continue
		}
		seen[result.DocumentID] = true
		sources = append(sources, result.DocumentID)
	}

	json.NewEncoder(w).Encode(askResponse{
		Answer:  reply,
		Sources: sources,
	})
}
