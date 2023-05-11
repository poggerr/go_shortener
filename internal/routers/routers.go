package routers

import (
	"github.com/go-chi/chi/v5"
	"github.com/poggerr/go_shortener/internal/handlers"
)

func Routers() chi.Router {
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Post("/", handlers.PostPage)
		r.Get("/{id}", handlers.GetPage)
	})
	return r
}
