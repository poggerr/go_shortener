package main

import (
	"bytes"
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/poggerr/go_shortener/internal/app"
	"github.com/poggerr/go_shortener/internal/handlers"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMainHendler(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
		oldUrl      string
		newUrl      string
		location    string
	}
	tests := []struct {
		name   string
		newUrl string
		oldUrl string
		want   want
	}{
		{
			name:   "testPost",
			oldUrl: "https://practicum.yandex.ru/",
			want: want{
				contentType: "text/plain ",
				statusCode:  201,
			},
		},
		{
			name:   "testGet",
			oldUrl: "https://practicum.yandex.ru/",
			want: want{
				contentType: "text/plain",
				statusCode:  307,
				location:    "https://practicum.yandex.ru/",
			},
		},
	}
	for _, tt := range tests {
		switch tt.name {
		case "testPost":
			t.Run(tt.name, func(t *testing.T) {
				body := []byte(tt.oldUrl)
				request := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
				w := httptest.NewRecorder()
				h := handlers.PostPage
				h(w, request)

				result := w.Result()

				assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))
				assert.Equal(t, tt.want.statusCode, result.StatusCode)
			})
		case "testGet":
			t.Run(tt.name, func(t *testing.T) {
				var target string
				for key, value := range app.MainMap {
					if value == tt.oldUrl {
						target = key
					}
				}
				request := httptest.NewRequest(http.MethodGet, "/{id}", nil)
				ctx := chi.NewRouteContext()
				ctx.URLParams.Add("id", target)
				request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, ctx))

				w := httptest.NewRecorder()
				h := handlers.GetPage
				h(w, request)

				result := w.Result()

				assert.Equal(t, tt.want.location, result.Header.Get("Location"))
				assert.Equal(t, tt.want.statusCode, result.StatusCode)
			})

		}

	}
}
