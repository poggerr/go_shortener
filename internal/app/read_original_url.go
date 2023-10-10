package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/poggerr/go_shortener/internal/logger"
)

// ReadOriginalURL хендлер получения оригинальной ссылки
func (a *App) ReadOriginalURL(res http.ResponseWriter, req *http.Request) {
	shortURL := chi.URLParam(req, "id")

	ans, err := a.storage.LongURL(shortURL)
	if err != nil {
		logger.Initialize().Info(err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	if a.storage.DB != nil {
		isDelete := a.storage.CheckLongURLIsDelete(shortURL)
		if isDelete {
			res.Header().Set("content-type", "text/plain ")
			res.WriteHeader(http.StatusGone)
			return
		}
	}

	res.Header().Set("content-type", "text/plain ")
	res.Header().Set("Location", ans)
	res.WriteHeader(307)

}
