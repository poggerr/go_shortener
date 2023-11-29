package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/poggerr/go_shortener/internal/models"
	mock_service "github.com/poggerr/go_shortener/internal/service/mocks"
	"github.com/poggerr/go_shortener/internal/utils"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const baseURL = "http://localhost:8080"

var errDumb = errors.New("dumb error")

func TestURLShortenerHandler_CreateShortURL(t *testing.T) {
	type fields struct {
		repo *mock_service.MockURLShortenerService
	}
	type want struct {
		status        int
		isURLInResult bool
	}
	tests := []struct {
		name    string
		body    string
		want    want
		prepare func(f *fields)
	}{
		{
			name: "valid link",
			body: "https://practicum.yandex.ru",
			want: want{
				status:        http.StatusCreated,
				isURLInResult: true,
			},
			prepare: func(f *fields) {
				gomock.InOrder(
					f.repo.EXPECT().Store(gomock.Any(), gomock.Any(), gomock.Any()).Return("1111", nil),
				)
			},
		},
		{
			name: "link already shortened",
			body: "https://practicum.yandex.ru",
			want: want{
				status:        http.StatusConflict,
				isURLInResult: true,
			},
			prepare: func(f *fields) {
				gomock.InOrder(
					f.repo.EXPECT().Store(gomock.Any(), gomock.Any(), gomock.Any()).Return("1111", errors.New("ссылка уже сокращена")),
				)
			},
		},
		{
			name: "invalid link",
			body: "yaru",
			want: want{
				status:        http.StatusBadRequest,
				isURLInResult: false,
			},
		},
		{
			name: "empty link",
			body: "",
			want: want{
				status:        http.StatusBadRequest,
				isURLInResult: false,
			},
		},
		{
			name: "repository error",
			body: "https://yandex.ru",
			want: want{
				status:        http.StatusInternalServerError,
				isURLInResult: false,
			},
			prepare: func(f *fields) {
				gomock.InOrder(
					f.repo.EXPECT().Store(gomock.Any(), gomock.Any(), gomock.Any()).Return("", errDumb),
				)
			},
		},
	}
	zerolog.SetGlobalLevel(zerolog.Disabled)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			mockRepo := mock_service.NewMockURLShortenerService(mockCtrl)

			f := fields{
				repo: mockRepo,
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}

			reader := strings.NewReader(tt.body)
			request := httptest.NewRequest(http.MethodPost, "/", reader)
			w := httptest.NewRecorder()
			h := NewURLShortener(baseURL, mockRepo)
			h.CreateShortURL(w, request)
			result := w.Result()

			urlResult, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, tt.want.status, result.StatusCode)
			assert.Equal(t, tt.want.isURLInResult, utils.IsURL(string(urlResult)))
		})
	}
}

type URLShortenResponse struct {
	Result string `json:"result"`
}

func TestURLShortener_CreateJSONShorten(t *testing.T) {
	type fields struct {
		repo *mock_service.MockURLShortenerService
	}
	type want struct {
		status int
		err    assert.ErrorAssertionFunc
	}
	tests := []struct {
		name    string
		body    string
		want    want
		prepare func(f *fields)
	}{
		{
			name: "valid link",
			body: `{"url":"https://practicum.yandex.ru"}`,
			want: want{
				status: http.StatusCreated,
				err:    assert.NoError,
			},
			prepare: func(f *fields) {
				gomock.InOrder(f.repo.EXPECT().Store(gomock.Any(), gomock.Any(), gomock.Any()).Return("1111", nil))
			},
		},
		{
			name: "link already shortened",
			body: `{"url":"https://practicum.yandex.ru"}`,
			want: want{
				status: http.StatusConflict,
				err:    assert.Error,
			},
			prepare: func(f *fields) {
				gomock.InOrder(f.repo.EXPECT().Store(gomock.Any(), gomock.Any(), gomock.Any()).Return("1111", errors.New("ссылка уже сокращена")))
			},
		},
		{
			name: "repository error",
			body: `{"url":"https://yandex.ru"}`,
			want: want{
				status: http.StatusInternalServerError,
				err:    assert.Error,
			},
			prepare: func(f *fields) {
				gomock.InOrder(
					f.repo.EXPECT().Store(gomock.Any(), gomock.Any(), gomock.Any()).Return("", errDumb),
				)
			},
		},
		{
			name: "invalid link",
			body: `{"url":"yaru"}`,
			want: want{
				status: http.StatusBadRequest,
				err:    assert.Error,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			mockRepo := mock_service.NewMockURLShortenerService(mockCtrl)

			f := fields{
				repo: mockRepo,
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}

			reader := strings.NewReader(tt.body)
			request := httptest.NewRequest(http.MethodPost, "/", reader)
			w := httptest.NewRecorder()
			h := NewURLShortener(baseURL, mockRepo)
			h.CreateJSONShorten(w, request)
			result := w.Result()

			require.Equal(t, tt.want.status, result.StatusCode)
			var resp URLShortenResponse
			err := json.NewDecoder(result.Body).Decode(&resp)
			if !tt.want.err(t, err, fmt.Sprintf("request: (%v)", tt.body)) {
				return
			}
			err = result.Body.Close()
			require.NoError(t, err)
		})
	}
}

func TestURLShortenerHandler_GetOriginalURL(t *testing.T) {
	type fields struct {
		repo *mock_service.MockURLShortenerService
	}
	type want struct {
		status   int
		location string
	}
	tests := []struct {
		name    string
		link    string
		want    want
		prepare func(f *fields)
	}{
		{
			name: "valid link",
			link: "http://localhost:8080",
			want: want{
				status:   http.StatusTemporaryRedirect,
				location: "https://practicum.yandex.ru",
			},
			prepare: func(f *fields) {
				gomock.InOrder(f.repo.EXPECT().Restore(gomock.Any(), gomock.Any()).Return("https://practicum.yandex.ru", nil))
			},
		},
		{
			name: "repository error",
			link: "http://localhost:8080",
			want: want{
				status:   http.StatusInternalServerError,
				location: "",
			},
			prepare: func(f *fields) {
				gomock.InOrder(f.repo.EXPECT().Restore(gomock.Any(), gomock.Any()).Return("", errDumb))
			},
		},
		{
			name: "deleted link",
			link: "http://localhost:8080",
			want: want{
				status:   http.StatusGone,
				location: "",
			},
			prepare: func(f *fields) {
				gomock.InOrder(f.repo.EXPECT().Restore(gomock.Any(), gomock.Any()).Return("", ErrLinkIsDeleted))
			},
		},
		{
			name: "invalid link",
			link: "http://localhost:8080",
			want: want{
				status:   http.StatusNoContent,
				location: "",
			},
			prepare: func(f *fields) {
				gomock.InOrder(f.repo.EXPECT().Restore(gomock.Any(), gomock.Any()).Return("", ErrLinkNotFound))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			mockRepo := mock_service.NewMockURLShortenerService(mockCtrl)

			f := fields{
				repo: mockRepo,
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}

			r := httptest.NewRequest(http.MethodGet, tt.link, nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", "1111")
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

			h := NewURLShortener(baseURL, mockRepo)
			w := httptest.NewRecorder()
			h.ReadOriginalURL(w, r)
			result := w.Result()
			err := result.Body.Close()
			require.NoError(t, err)

			require.Equal(t, tt.want.status, result.StatusCode)
			assert.Equal(t, tt.want.location, result.Header.Get("Location"))
		})
	}
}

func TestURLShortener_HandleDelete(t *testing.T) {
	type fields struct {
		repo *mock_service.MockURLShortenerService
	}
	type want struct {
		status int
	}
	tests := []struct {
		name    string
		reqBody string
		want    want
		prepare func(f *fields)
	}{
		{
			name:    "valid request",
			reqBody: `["111","222"]`,
			want: want{
				status: http.StatusAccepted,
			},
			prepare: func(f *fields) {
				gomock.InOrder(
					f.repo.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes(),
				)
			},
		},
		{
			name:    "invalid request",
			reqBody: `["111","222]`,
			want: want{
				status: http.StatusInternalServerError,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			mockRepo := mock_service.NewMockURLShortenerService(mockCtrl)
			f := fields{
				repo: mockRepo,
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}

			reader := strings.NewReader(tt.reqBody)
			request := httptest.NewRequest(http.MethodDelete, "/", reader)
			w := httptest.NewRecorder()
			h := NewURLShortener(baseURL, mockRepo)
			h.DeleteUrls(w, request)
			result := w.Result()
			err := result.Body.Close()

			require.NoError(t, err)
			require.Equal(t, tt.want.status, result.StatusCode)
		})
	}
}

func TestURLShortener_GetUrlsByUser(t *testing.T) {
	type fields struct {
		repo *mock_service.MockURLShortenerService
	}
	type want struct {
		status int
		result string
	}
	tests := []struct {
		name    string
		reqBody string
		want    want
		prepare func(f *fields)
	}{
		{
			name: "single item",
			want: want{
				status: http.StatusOK,
				result: `[
					  {
					    "short_url": "http://localhost:8080/1111",
					    "original_url": "https://practicum.yandex.ru"
					  }
					]`,
			},
			prepare: func(f *fields) {
				gomock.InOrder(
					f.repo.EXPECT().GetUserStorage(gomock.Any(), gomock.Any(), gomock.Any()).
						Return(map[string]string{"1111": "https://practicum.yandex.ru"}, nil),
				)
			},
		},
		// TODO переписать
		//{
		//	name: "couple item",
		//	want: want{
		//		status: http.StatusOK,
		//		result: `[
		//			  {
		//			    "short_url": "http://localhost:8080/1111",
		//			    "original_url": "https://practicum.yandex.ru"
		//			  },
		//			  {
		//			    "short_url": "http://localhost:8080/2222",
		//			    "original_url": "https://yandex.ru"
		//			  }
		//			]`,
		//	},
		//	prepare: func(f *fields) {
		//		gomock.InOrder(
		//			f.repo.EXPECT().GetUserStorage(gomock.Any(), gomock.Any(), gomock.Any()).
		//				Return(map[string]string{"1111": "https://practicum.yandex.ru", "2222": "https://yandex.ru"}, nil),
		//		)
		//	},
		//},
		{
			name: "empty",
			want: want{
				status: http.StatusNoContent,
				result: "",
			},
			prepare: func(f *fields) {
				gomock.InOrder(
					f.repo.EXPECT().GetUserStorage(gomock.Any(), gomock.Any(), gomock.Any()).
						Return(nil, errors.New("some error")),
				)
			},
		},
	}
	zerolog.SetGlobalLevel(zerolog.Disabled)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			mockRepo := mock_service.NewMockURLShortenerService(mockCtrl)

			f := fields{
				repo: mockRepo,
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}

			request := httptest.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()
			h := NewURLShortener(baseURL, mockRepo)
			h.GetUrlsByUser(w, request)
			result := w.Result()
			assert.Equal(t, tt.want.status, result.StatusCode)

			buf := new(bytes.Buffer)
			_, err := buf.ReadFrom(result.Body)
			require.NoError(t, err)
			if tt.want.result == "" || buf.String() == "" {
				assert.EqualValues(t, tt.want.result, buf.String())
			} else {
				assert.JSONEq(t, tt.want.result, buf.String())
			}
			err = result.Body.Close()
			require.NoError(t, err)
		})
	}
}

var singleItem = models.BatchList{
	{
		CorrelationID: "xxxx",
		OriginalURL:   "https://ya.ru",
		ShortURL:      "http://localhost:8080/1111",
	},
}

var coupleItem = models.BatchList{
	{
		CorrelationID: "xxxx",
		OriginalURL:   "https://ya.ru",
		ShortURL:      "http://localhost:8080/1111",
	},
	{
		CorrelationID: "yyyy",
		OriginalURL:   "https://yandex.ru",
		ShortURL:      "http://localhost:8080/2222",
	},
}

func TestURLShortener_CreateBatch(t *testing.T) {
	type fields struct {
		repo *mock_service.MockURLShortenerService
	}
	type want struct {
		status int
		result string
	}
	tests := []struct {
		name    string
		reqBody string
		want    want
		prepare func(f *fields)
	}{
		{
			name:    "empty body",
			reqBody: "",
			want:    want{status: http.StatusBadRequest, result: ""},
			prepare: nil,
		},
		{
			name: "single item",
			reqBody: `
					[
					   {
						 "correlation_id": "xxxx",
						 "original_url": "https://ya.ru"
					   }
					]
			`,
			want: want{
				status: http.StatusCreated,
				result: `
						[
						   {
							 "correlation_id": "xxxx",
							 "short_url": "http://localhost:8080/1111",
                             "original_url": "https://ya.ru"
						   }
						]
				`,
			},
			prepare: func(f *fields) {
				gomock.InOrder(
					f.repo.EXPECT().StoreBatch(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
						Return(singleItem, nil),
				)
			},
		},
		{
			name: "couple items",
			reqBody: `
				[
				  {
					"correlation_id": "xxxx",
					"original_url": "https://ya.ru"
				  },
				  {
					"correlation_id": "yyyy",
					"original_url": "https://yandex.ru"
				  }
				]
			`,
			want: want{
				status: http.StatusCreated,
				result: `
					[
					   {
						 "correlation_id": "xxxx",
						 "short_url": "http://localhost:8080/1111",
                         "original_url": "https://ya.ru"
					   },
					   {
						 "correlation_id": "yyyy",
						 "short_url": "http://localhost:8080/2222",
						 "original_url": "https://yandex.ru"
					   }
					 ]
				`,
			},
			prepare: func(f *fields) {
				gomock.InOrder(
					f.repo.EXPECT().StoreBatch(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
						Return(coupleItem, nil),
				)
			},
		},
		{
			name: "internal server error",
			reqBody: `
				[
				  {
					"correlation_id": "xxxx",
					"original_url": "https://ya.ru"
				  }
				]
			`,
			want: want{
				status: http.StatusInternalServerError,
			},
			prepare: func(f *fields) {
				gomock.InOrder(
					f.repo.EXPECT().StoreBatch(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
						Return(nil, errDumb),
				)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			mockRepo := mock_service.NewMockURLShortenerService(mockCtrl)

			f := fields{
				repo: mockRepo,
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}

			reader := strings.NewReader(tt.reqBody)
			request := httptest.NewRequest(http.MethodPost, "/", reader)
			w := httptest.NewRecorder()
			h := NewURLShortener(baseURL, mockRepo)
			h.CreateBatch(w, request)
			result := w.Result()
			assert.Equal(t, tt.want.status, result.StatusCode)

			buf := new(bytes.Buffer)
			_, err := buf.ReadFrom(result.Body)
			require.NoError(t, err)
			if tt.want.result == "" || buf.String() == "" ||
				tt.want.status == http.StatusBadRequest ||
				tt.want.status == http.StatusInternalServerError {
				assert.EqualValues(t, tt.want.result, buf.String())
			} else {
				assert.JSONEq(t, tt.want.result, buf.String())
			}
			err = result.Body.Close()
			require.NoError(t, err)
		})
	}
}

func TestURLShortener_Ping(t *testing.T) {
	type fields struct {
		repo *mock_service.MockURLShortenerService
	}
	type want struct {
		status int
	}
	tests := []struct {
		name    string
		want    want
		prepare func(f *fields)
	}{
		{
			name: "ok",
			want: want{status: http.StatusOK},
			prepare: func(f *fields) {
				gomock.InOrder(
					f.repo.EXPECT().Ping(gomock.Any()).Return(nil),
				)
			},
		},
		{
			name: "not ok",
			want: want{status: http.StatusInternalServerError},
			prepare: func(f *fields) {
				gomock.InOrder(
					f.repo.EXPECT().Ping(gomock.Any()).Return(errDumb),
				)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			mockRepo := mock_service.NewMockURLShortenerService(mockCtrl)

			f := fields{
				repo: mockRepo,
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}

			request := httptest.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()
			h := NewURLShortener(baseURL, mockRepo)
			h.DBConnect(w, request)
			result := w.Result()
			assert.Equal(t, tt.want.status, result.StatusCode)
			err := result.Body.Close()
			require.NoError(t, err)
		})
	}
}
