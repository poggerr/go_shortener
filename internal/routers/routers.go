package routers

import (
	"github.com/go-chi/chi/v5"
	"github.com/poggerr/go_shortener/internal/app"
	"github.com/poggerr/go_shortener/internal/app/storage"
	"github.com/poggerr/go_shortener/internal/config"
	"github.com/poggerr/go_shortener/internal/gzip"
	"github.com/poggerr/go_shortener/internal/logger"
)

func Router(cfg *config.Config, strg *storage.Storage) chi.Router {
	r := chi.NewRouter()
	newApp := app.NewApp(cfg, strg)
	r.Use(gzip.GzipMiddleware)
	r.Use(logger.WithLoggingReq)
	r.Use(logger.WithLoggingRes)
	r.Post("/", newApp.CreateShortUrl)
	r.Post("/api/shorten", newApp.CreateJsonShorten)
	r.Get("/{id}", newApp.ReadOldUrl)
	return r
}
