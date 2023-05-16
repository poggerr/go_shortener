package routers

import (
	"bytes"
	"github.com/poggerr/go_shortener/internal/app/storage"
	"github.com/poggerr/go_shortener/internal/config"
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

func TestHandlersPost(t *testing.T) {
	ts := httptest.NewServer(Router(&cfg, strg))
	defer ts.Close()

	var testTable = []struct {
		method      string
		url         string
		contentType string
		status      int
		location    string
	}{
		{method: "POST", url: "https://practicum.yandex.ru/", contentType: "text/plain", status: 201},
		{method: "POST", url: "https://www.google.com/", contentType: "text/plain", status: 201},
		{method: "POST", url: "", contentType: "text/plain; charset=utf-8", status: 400},
		{method: "GET", url: "https://practicum.yandex.ru/", status: 200, location: "https://practicum.yandex.ru/"},
		{method: "GET", url: "https://www.google.com/", status: 200, location: "https://www.google.com/"},
	}

	for _, v := range testTable {
		switch v.method {
		case "POST":
			resp, respBody := testRequestPost(t, ts, v.method, "/", v.url)
			assert.Equal(t, v.status, resp.StatusCode)
			assert.Equal(t, v.contentType, resp.Header.Get("Content-Type"))
			if v.url != "" {
				m := strings.Split(respBody, "/")

				mainMap[v.url] = m[3]
			}
		case "GET":
			newUrl := "/"
			for key, value := range mainMap {
				if key == v.location {
					newUrl += value
				}
			}
			resp, _ := testRequestPost(t, ts, http.MethodGet, newUrl, "")
			assert.Equal(t, v.status, resp.StatusCode)
			assert.Equal(t, v.contentType, resp.Header.Get("Location"))
		}

	}

}
