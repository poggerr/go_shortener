package app

import (
	"encoding/json"
	"github.com/poggerr/go_shortener/internal/authorization"
	"github.com/poggerr/go_shortener/internal/logger"
	"github.com/poggerr/go_shortener/internal/models"
	"github.com/poggerr/go_shortener/internal/service_create_short_url"
	"io"
	"net/http"
)

// CreateJSONShorten хендлер создания ссылки из json
func (a *App) CreateJSONShorten(res http.ResponseWriter, req *http.Request) {
	userID := authorization.FromContext(req.Context())

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return
	}

	var url models.URL

	err = json.Unmarshal(body, &url)
	if err != nil {
		logger.Initialize().Info(err)
	}

	shortURL := service_create_short_url.CreateShortURL(url.LongURL)
	a.storage.Save(shortURL, url.LongURL)

	var short string

	switch {
	case a.storage.DB != nil:
		short, err = a.storage.SaveToDB(string(body), shortURL, userID)
		if err != nil {
			shortenMap := make(map[string]string)

			shortenMap["result"] = shortURL

			marshal, err := json.Marshal(shortenMap)
			if err != nil {
				logger.Initialize().Info(err)
			}

			res.Header().Set("content-type", "application/json ")
			res.WriteHeader(http.StatusConflict)
			res.Write(marshal)
			return
		}
	}

	shortURL = a.cfg.DefURL + "/" + short

	shortenMap := make(map[string]string)

	shortenMap["result"] = shortURL

	marshal, err := json.Marshal(shortenMap)
	if err != nil {
		logger.Initialize().Info(err)
	}

	res.Header().Set("content-type", "application/json ")
	res.WriteHeader(http.StatusCreated)

	res.Write(marshal)

}
