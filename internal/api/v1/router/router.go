package router

import (
	"net/http"

	"github.com/harishnagaraju/astramind/internal/api/v1/handlers"
	"github.com/harishnagaraju/astramind/internal/api/v1/request"
	"github.com/harishnagaraju/astramind/internal/features/kb"
	"github.com/harishnagaraju/astramind/internal/infrastructure/ai"
)

type Config struct {
	ProviderName    string
	Model           string
	Version         string
	APIKey          string
	ProviderManager *ai.ProviderManager
	KnowledgeBase   *kb.Manager
}

func New(config Config) http.Handler {
	handler := handlers.New(handlers.Config{
		ProviderName: config.ProviderName,
		Model: config.Model,
		Version: config.Version,
		APIKey: config.APIKey,
		ProviderManager: config.ProviderManager,
		KnowledgeBase: config.KnowledgeBase,
	})

	mux := http.NewServeMux()
	mux.HandleFunc("/health", handler.Health)
	mux.HandleFunc("/status", handler.Status)
	mux.HandleFunc("/version", handler.Version)
	mux.HandleFunc("/providers", handler.Providers)
	mux.HandleFunc("/models", handler.Models)
	mux.HandleFunc("/chat", handler.Chat)
	mux.HandleFunc("/chat/stream", handler.ChatStream)
	mux.HandleFunc("/documents", handler.Documents)
	mux.HandleFunc("/documents/", func(w http.ResponseWriter, r *http.Request) {
		if len(r.URL.Path) >= len("/documents/") && len(r.URL.Path) > len("/documents/") && len(r.URL.Path)-len("/documents/") > len("/analyze") && r.URL.Path[len(r.URL.Path)-len("/analyze"):] == "/analyze" {
			handler.AnalyzeDocument(w, r)
			return
		}
		handler.Document(w, r)
	})
	mux.HandleFunc("/embeddings", handler.Embeddings)
	mux.HandleFunc("/search", handler.Search)

	return request.Middleware(mux)
}
