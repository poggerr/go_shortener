package middlewares

import (
	"github.com/poggerr/go_shortener/internal/logger"
	"net/http"
	"time"
)

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
