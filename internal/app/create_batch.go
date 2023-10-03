package app

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/poggerr/go_shortener/internal/logger"
	"github.com/poggerr/go_shortener/internal/models"
)

// CreateBatch Хендлер создания нескольких ссылок
func (a *App) CreateBatch(res http.ResponseWriter, req *http.Request) {
	var list models.BatchList
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &list)
	if err != nil {
		logger.Initialize().Info(err)
	}

	list = a.storage.SaveMultipleToDB(list, a.cfg.DefURL)

	marshal, err := json.Marshal(list)
	if err != nil {
		logger.Initialize().Info("Ошибка при формировании ответа ", err)
	}

	res.Header().Set("content-type", "application/json ")
	res.WriteHeader(http.StatusCreated)
	res.Write(marshal)

}
