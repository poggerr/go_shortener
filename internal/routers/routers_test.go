package routers

import (
	"bytes"
	"encoding/json"
	"github.com/poggerr/go_shortener/internal/app/storage"
	"github.com/poggerr/go_shortener/internal/config"
	"github.com/poggerr/go_shortener/internal/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var mainMap = make(map[string]string)
var cfg = config.NewDefConf()
var strg = storage.NewStorage()

func testRequestPost(t *testing.T, ts *httptest.Server, method,
	path string, oldUrl string) (*http.Response, string) {

	req, err := http.NewRequest(method, ts.URL+path, bytes.NewBuffer([]byte(oldUrl)))
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

func testRequestJson(t *testing.T, ts *httptest.Server, method, path string, longUrl string) (*http.Response, string) {
	longUrlMap := make(map[string]string)
	longUrlMap["url"] = longUrl
	marshal, _ := json.Marshal(longUrlMap)

	req, err := http.NewRequest(method, ts.URL+path, bytes.NewBuffer(marshal))
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

func TestHandlersPost(t *testing.T) {
	logger.Initialize()
	ts := httptest.NewServer(Router(&cfg, strg))
	defer ts.Close()

	var testTable = []struct {
		api         string
		method      string
		url         string
		contentType string
		status      int
		location    string
	}{
		{api: "/", method: "POST", url: "https://practicum.yandex.ru/", contentType: "text/plain", status: 201},
		{api: "/", method: "POST", url: "https://www.google.com/", contentType: "text/plain", status: 201},
		{api: "/", method: "POST", url: "", contentType: "text/plain; charset=utf-8", status: 400},
		{api: "/id", method: "GET", url: "https://practicum.yandex.ru/", status: 200, location: "https://practicum.yandex.ru/"},
		{api: "/id", method: "GET", url: "https://www.google.com/", status: 200, location: "https://www.google.com/"},
		{api: "/api/shorten", method: "POST", url: "https://practicum.yandex.ru/", contentType: "application/json", status: 200},
		{api: "/api/shorten", method: "POST", url: "https://www.google.com/", contentType: "application/json", status: 200},
	}

	for _, v := range testTable {
		switch v.api {
		case "/":
			resp, respBody := testRequestPost(t, ts, v.method, v.api, v.url)
			assert.Equal(t, v.status, resp.StatusCode)
			assert.Equal(t, v.contentType, resp.Header.Get("Content-Type"))
			if v.url != "" {
				m := strings.Split(respBody, "/")

				mainMap[v.url] = m[3]
			}
		case "/id":
			newUrl := "/"
			for key, value := range mainMap {
				if key == v.location {
					newUrl += value
				}
			}
			resp, _ := testRequestPost(t, ts, http.MethodGet, newUrl, "")
			assert.Equal(t, v.status, resp.StatusCode)
			assert.Equal(t, v.contentType, resp.Header.Get("Location"))
		case "/api/shorten":
			resp, _ := testRequestJson(t, ts, v.method, v.api, v.url)
			assert.Equal(t, v.status, resp.StatusCode)
			assert.Equal(t, v.contentType, resp.Header.Get("Content-Type"))
		}

	}

}
