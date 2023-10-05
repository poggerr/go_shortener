package app

import (
	"github.com/poggerr/go_shortener/internal/authorization"
	"github.com/poggerr/go_shortener/internal/logger"
	"github.com/poggerr/go_shortener/internal/servicecreateshorturl"
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

	shortURL := servicecreateshorturl.CreateShortURL(string(body))
	a.storage.Save(shortURL, string(body))

	switch {
	case a.storage.DB != nil:
		short, err := a.storage.SaveToDB(string(body), shortURL, userID)
		if err != nil {
			logger.Initialize().Info(err)
			res.Header().Set("content-type", "text/plain; charset=utf-8")
			res.WriteHeader(http.StatusConflict)
			shortURL = a.cfg.DefURL + "/" + short
			res.Write([]byte(shortURL))
			return
		}
	}

	shortURL = a.cfg.DefURL + "/" + shortURL

	res.Header().Set("content-type", "text/plain; charset=utf-8")

	res.WriteHeader(http.StatusCreated)

	res.Write([]byte(shortURL))

}
