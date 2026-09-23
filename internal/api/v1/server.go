package v1

import (
	"net/http"

	"github.com/harishnagaraju/astramind/internal/api/v1/router"
)

func New(config Config) http.Handler {
	return router.New(router.Config{
		ProviderName: config.ProviderName,
		Model: config.Model,
		Version: config.Version,
		APIKey: config.APIKey,
		ProviderManager: config.ProviderManager,
		KnowledgeBase: config.KnowledgeBase,
	})
}

func Mount(mux *http.ServeMux, config Config) {
	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", New(config)))
}
