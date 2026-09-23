package v1

import (
	"net/http"

	"github.com/harishnagaraju/astramind/internal/api/v1/router"
)

func New(config Config) http.Handler {
	return router.New(router.Config{
		ProviderName: config.ProviderName,
		Model:        config.Model,
		Version:      config.Version,
	})
}

func Mount(mux *http.ServeMux, config Config) {
	handler := New(config)
	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", handler))
}
