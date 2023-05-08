package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/poggerr/go_shortener/config/flags"
	"github.com/poggerr/go_shortener/internal/app"
	"io"
	"net/http"
)

func GetPage(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(res, "Only GET requests are allowed!", http.StatusBadRequest)
		return
	}
	id := chi.URLParam(req, "id")
	ans := app.UnShorting(id)
	res.Header().Set("Location", ans)
	res.WriteHeader(307)

}

func PostPage(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Only Post requests are allowed!", http.StatusBadRequest)
		return
	}

	if err := req.ParseForm(); err != nil {
		res.Write([]byte(err.Error()))
		return
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return
	}

	local := flags.DefUrl + "/"

	local += app.Shorting(string(body))

	res.Header().Set("content-type", "text/plain ")

	res.WriteHeader(http.StatusCreated)

	res.Write([]byte(local))

}
