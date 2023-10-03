package app

import (
	"context"
	"net/http"
	"time"

	"github.com/poggerr/go_shortener/internal/logger"
)

// DBConnect хендлер проверки подключения к БД
func (a *App) DBConnect(res http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := a.db.PingContext(ctx); err != nil {
		logger.Initialize().Error("Ошибка при подключении к БД ", err)
		res.WriteHeader(http.StatusInternalServerError)
	}
	res.WriteHeader(http.StatusOK)
}
