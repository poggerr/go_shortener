// Package middlewares содержит все основные middleware проекта
package middlewares

import (
	"fmt"
	"net/http"
	"time"

	"github.com/poggerr/go_shortener/internal/logger"
)

// WithLogging Логирование каждого запроса
func WithLogging(h http.Handler) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		responseData := &logger.ResponseData{}
		lw := logger.LoggingResponseWriter{
			ResponseWriter: w,
			ResponseData:   responseData,
		}

		h.ServeHTTP(&lw, r)

		duration := time.Since(time.Now())

		logger.Log.Infoln(
			"uri", r.RequestURI,
			"method", r.Method,
			"status", responseData.Status,
			"duration", duration,
			"size", responseData.Size,
		)
	}
	return http.HandlerFunc(logFn)
}

func WithCookies(h http.Handler) http.Handler {
	fn := func(res http.ResponseWriter, req *http.Request) {
		c, err := req.Cookie("session_token")
		if err != nil {
			logger.Initialize().Info("Ошибка при получении Cookie ", err)
		}

		fmt.Println(c)

	}
	return http.HandlerFunc(fn)
}
