package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog"
	"github.com/poggerr/go_shortener/internal/authorization"
	"github.com/poggerr/go_shortener/internal/gzip"
	"github.com/poggerr/go_shortener/internal/handlers"
	"github.com/poggerr/go_shortener/internal/service"
	"github.com/rs/zerolog/log"
	"net/http"
)

// Server возвращает сервер
func Server(addr, baseURL string, repo service.URLShortenerService) *http.Server {
	Repo := repo
	handler := handlers.NewURLShortener(baseURL, Repo)
	r := chi.NewRouter()
	r.Use(httplog.RequestLogger(log.Logger))

	r.Use(authorization.AuthMiddleware, gzip.GzipMiddleware)
	r.Post("/", handler.CreateShortURL)
	r.Post("/api/shorten", handler.CreateJSONShorten)
	r.Get("/{id}", handler.ReadOriginalURL)
	r.Get("/ping", handler.DBConnect)
	r.Post("/api/shorten/batch", handler.CreateBatch)
	r.Get("/api/user/urls", handler.GetUrlsByUser)
	r.Delete("/api/user/urls", handler.DeleteUrls)

	server := &http.Server{
		Addr:           addr,
		Handler:        r,
		TLSConfig:      nil,
		MaxHeaderBytes: 16 * 1024,
	}

	return server
}
