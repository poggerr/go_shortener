package app

import (
	"github.com/poggerr/go_shortener/internal/authorization"
	"github.com/poggerr/go_shortener/internal/logger"
	"github.com/poggerr/go_shortener/internal/service_create_short_url"
	"io"
	"net/http"
)

// CreateShortURL хендлер создания короткой ссылки
func (a *App) CreateShortURL(res http.ResponseWriter, req *http.Request) {
	userID := authorization.FromContext(req.Context())

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return
	}

	shortURL := service_create_short_url.CreateShortURL(string(body))
	a.storage.Save(shortURL, string(body))
	if a.storage.DB != nil {
		a.storage.SaveToDB(string(body), shortURL, userID)
	}

	shortURL = a.cfg.DefURL + "/" + shortURL
	if err != nil {
		logger.Initialize().Info(err)
		res.Header().Set("content-type", "text/plain; charset=utf-8")
		res.WriteHeader(http.StatusConflict)
		res.Write([]byte(shortURL))
		return
	}

	res.Header().Set("content-type", "text/plain; charset=utf-8")

	res.WriteHeader(http.StatusCreated)

	res.Write([]byte(shortURL))

}
