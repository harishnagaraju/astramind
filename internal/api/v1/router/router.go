package router

import (
	"net/http"

	"github.com/harishnagaraju/astramind/internal/api/v1/handlers"
	"github.com/harishnagaraju/astramind/internal/api/v1/request"
)

type Config struct {
	ProviderName string
	Model        string
	Version      string
}

func New(config Config) http.Handler {
	handler := handlers.New(handlers.Config{
		ProviderName: config.ProviderName,
		Model:        config.Model,
		Version:      config.Version,
	})

	mux := http.NewServeMux()
	mux.HandleFunc("/health", handler.Health)
	mux.HandleFunc("/status", handler.Status)
	mux.HandleFunc("/version", handler.Version)

	return request.Middleware(mux)
}
