package handlers

import (
	"github.com/golang/mock/gomock"
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

const baseURL = "http://localhost:8080/"

func TestURLShortenerHandler_HandlePostShortenPlain(t *testing.T) {
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
			body: "https://habr.com/ru/post/66931/",
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
		//{
		//	name: "link already shortened",
		//	body: "https://habr.com/ru/post/66931/",
		//	want: want{
		//		status:        http.StatusConflict,
		//		isURLInResult: true,
		//	},
		//	prepare: func(f *fields) {
		//		gomock.InOrder(
		//			f.repo.EXPECT().Store(gomock.Any(), gomock.Any(), gomock.Any()).Return("ссылка уже сокращена", nil),
		//		)
		//	},
		//},
		//{
		//	name: "repository error",
		//	body: "https://habr.com/ru/post/66931/",
		//	want: want{
		//		status:        http.StatusInternalServerError,
		//		isURLInResult: false,
		//	},
		//	prepare: func(f *fields) {
		//		gomock.InOrder(
		//			f.repo.EXPECT().Store(gomock.Any(), gomock.Any(), gomock.Any()).Return("", nil),
		//		)
		//	},
		//},
		//{
		//	name: "invalid link",
		//	body: "yaru",
		//	want: want{
		//		status:        http.StatusBadRequest,
		//		isURLInResult: false,
		//	},
		//},
		//{
		//	name: "empty link",
		//	body: "",
		//	want: want{
		//		status:        http.StatusBadRequest,
		//		isURLInResult: false,
		//	},
		//},
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
