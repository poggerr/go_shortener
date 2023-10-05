package routers

import (
	"database/sql"
	"github.com/poggerr/go_shortener/internal/async"
	"github.com/poggerr/go_shortener/internal/authorization"
	"github.com/poggerr/go_shortener/internal/middlewares"
	"github.com/poggerr/go_shortener/internal/storage"
	"net/http"
	"net/http/pprof"

	"github.com/go-chi/chi/v5"
	"github.com/poggerr/go_shortener/internal/app"
	"github.com/poggerr/go_shortener/internal/config"
)

func Router(cfg *config.Config, strg *storage.Storage, db *sql.DB, repo *async.URLRepo) chi.Router {
	r := chi.NewRouter()
	newApp := app.NewApp(cfg, strg, db, repo)
	r.Use(middlewares.WithLogging, authorization.AuthMiddleware)
	r.Post("/", newApp.CreateShortURL)
	r.Post("/api/shorten", newApp.CreateJSONShorten)
	r.Get("/{id}", newApp.ReadOriginalURL)
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
