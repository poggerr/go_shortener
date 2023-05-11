package routers

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

//var mainMap = make(map[string]string)

func testRequest(t *testing.T, ts *httptest.Server, method,
	path string, oldUrl string) (*http.Response, string) {
	reqBody := []byte(oldUrl)
	req, err := http.NewRequest(method, ts.URL+path, bytes.NewBuffer(reqBody))
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

func TestHandlersPost(t *testing.T) {
	ts := httptest.NewServer(Routers())
	defer ts.Close()

	var testTable = []struct {
		url         string
		contentType string
		status      int
	}{
		{"https://practicum.yandex.ru/", "text/plain ", 201},
	}

	for _, v := range testTable {
		resp, _ := testRequest(t, ts, "POST", "/", v.url)
		assert.Equal(t, v.status, resp.StatusCode)
		assert.Equal(t, v.contentType, resp.Header.Get("Content-Type"))
	}

	//t.Run(testPost.name, func(t *testing.T) {
	//	req := resty.New().R()
	//	req.Method = http.MethodPost
	//	req.URL = srv.URL
	//	req.Body = []byte(testPost.oldUrl)
	//
	//	resp, err := req.Send()
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//	fmt.Println(resp.Body())
	//
	//	assert.Equal(t, testPost.want.contentType, resp.Header().Get("Content-Type"))
	//	assert.Equal(t, testPost.want.statusCode, resp.StatusCode())
	//})

	//for _, tt := range tests {
	//	switch tt.name {
	//	case "testPost":
	//		t.Run(tt.name, func(t *testing.T) {
	//			body := []byte(tt.oldUrl)
	//			request := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
	//			w := httptest.NewRecorder()
	//			h := PostPage
	//			h(w, request)
	//
	//			result := w.Result()
	//
	//			read, err := io.ReadAll(w.Body)
	//			if err != nil {
	//				log.Printf("Error ReadAll %#v\n", err)
	//			}
	//
	//			newUrl = string(read)
	//
	//			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))
	//			assert.Equal(t, tt.want.statusCode, result.StatusCode)
	//		})
	//	case "testGet":
	//		t.Run(tt.name, func(t *testing.T) {
	//
	//			request := httptest.NewRequest(http.MethodGet, "/{id}", nil)
	//			ctx := chi.NewRouteContext()
	//			ctx.URLParams.Add("id", newUrl)
	//			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, ctx))
	//
	//			w := httptest.NewRecorder()
	//			h := GetPage
	//			h(w, request)
	//
	//			result := w.Result()
	//
	//			assert.Equal(t, tt.want.location, result.Header.Get("Location"))
	//			assert.Equal(t, tt.want.statusCode, result.StatusCode)
	//		})
	//
	//	}
	//
	//}
}

//tests := []struct {
//	name   string
//	newUrl string
//	oldUrl string
//	want   want
//}{
//	{
//		name:   "testPost",
//		oldUrl: "https://practicum.yandex.ru/",
//		want: want{
//			contentType: "text/plain ",
//			statusCode:  201,
//		},
//	},
//	{
//		name:   "testGet",
//		oldUrl: "https://practicum.yandex.ru/",
//		want: want{
//			contentType: "text/plain",
//			statusCode:  307,
//			location:    "https://practicum.yandex.ru/",
//		},
//	},
//}
