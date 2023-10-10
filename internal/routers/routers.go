package routers

import (
	"database/sql"
	"github.com/go-chi/chi/v5"
	"github.com/poggerr/go_shortener/internal/app"
	"github.com/poggerr/go_shortener/internal/app/middlewares"
	"github.com/poggerr/go_shortener/internal/app/service"
	"github.com/poggerr/go_shortener/internal/app/storage"
	"github.com/poggerr/go_shortener/internal/config"
	"github.com/poggerr/go_shortener/internal/gzip"
	"net/http"
	"net/http/pprof"
)

func Router(cfg *config.Config, strg *storage.Storage, db *sql.DB, repo *service.URLRepo) chi.Router {
	r := chi.NewRouter()
	newApp := app.NewApp(cfg, strg, db, repo)
	r.Use(middlewares.WithLogging, gzip.GzipMiddleware)
	r.Post("/", newApp.CreateShortURL)
	r.Post("/api/shorten", newApp.CreateJSONShorten)
	r.Get("/{id}", newApp.ReadOldURL)
	r.Get("/ping", newApp.DBConnect)
	r.Post("/api/shorten/batch", newApp.CreateBatch)
	r.Get("/api/user/urls", newApp.GetUrlsByUser)
	r.Delete("/api/user/urls", newApp.DeleteUrls)

	r.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
	r.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
	r.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
	r.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	r.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))
	r.Handle("/debug/pprof/{cmd}", http.HandlerFunc(pprof.Index))
	return r
}
