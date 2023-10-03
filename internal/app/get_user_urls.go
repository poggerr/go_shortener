package app

import (
	"encoding/json"
	"net/http"

	"github.com/poggerr/go_shortener/internal/authorization"
	"github.com/poggerr/go_shortener/internal/logger"
)

// GetUrlsByUser хендлер получения списка ссылок пользователя
func (a *App) GetUrlsByUser(res http.ResponseWriter, req *http.Request) {
	userID := authorization.FromContext(req.Context())

	strg, err := a.storage.GetUrlsByUserID(userID, a.cfg.DefURL)
	if err != nil {
		logger.Initialize().Info(err)
		res.WriteHeader(http.StatusUnauthorized)
		return
	}

	marshal, err := json.Marshal(strg)
	if err != nil {
		logger.Initialize().Info(err)
		return
	}

	res.Header().Set("content-type", "application/json ")
	res.WriteHeader(http.StatusOK)
	res.Write(marshal)

}
