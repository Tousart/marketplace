package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/tousart/marketplace/config"
)

func CreateAndRunServer(r chi.Router, cfg *config.Config) error {
	httpServer := &http.Server{
		Addr:    cfg.HTTP.Address,
		Handler: r,
	}

	return httpServer.ListenAndServe()
}
