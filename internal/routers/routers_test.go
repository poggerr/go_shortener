package routers

import (
	"bytes"
	"compress/gzip"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/poggerr/go_shortener/internal/app/storage"
	"github.com/poggerr/go_shortener/internal/config"
	"github.com/poggerr/go_shortener/internal/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

var mainMap = make(map[string]string)

func NewDefConf() config.Config {
	return config.Config{
		Serv:   ":8080",
		DefURL: "http://localhost:8080",
		Path:   "/tmp/short-url-db3.json",
		DB:     "host=localhost user=username password=userpassword dbname=shortener sslmode=disable",
	}
}

var cfg = NewDefConf()
var strg = storage.NewStorage("/tmp/short-url-db.json", connectDB())

func connectDB() *sql.DB {
	db, err := sql.Open("pgx", cfg.DB)
	if err != nil {
		logger.Initialize().Error("Ошибка при подключении к БД ", err)
	}
	defer db.Close()
	return db
}

func testRequestPost(t *testing.T, ts *httptest.Server, method,
	path string, oldURL string) (*http.Response, string) {

	req, err := http.NewRequest(method, ts.URL+path, bytes.NewBuffer([]byte(oldURL)))
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

func testRequestJSON(t *testing.T, ts *httptest.Server, method, path string, longURL string) (*http.Response, string) {
	longURLMap := make(map[string]string)
	longURLMap["url"] = longURL
	marshal, _ := json.Marshal(longURLMap)

	req, err := http.NewRequest(method, ts.URL+path, bytes.NewBuffer(marshal))
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

func TestHandlersPost(t *testing.T) {
	logger.Initialize()
	ts := httptest.NewServer(Router(&cfg, strg, strg.DB))
	defer ts.Close()

	var testTable = []struct {
		api         string
		method      string
		url         string
		contentType string
		status      int
		location    string
	}{
		{api: "/", method: "POST", url: "https://prabicum.yandex.ru/", contentType: "text/plain; charset=utf-8", status: 409},
		{api: "/", method: "POST", url: "https://www.gjle.com/", contentType: "text/plain; charset=utf-8", status: 409},
	}

	for _, v := range testTable {
		switch v.api {
		case "/":
			resp, _ := testRequestPost(t, ts, v.method, v.api, v.url)
			defer resp.Body.Close()
			assert.Equal(t, v.status, resp.StatusCode)
			assert.Equal(t, v.contentType, resp.Header.Get("Content-Type"))
		case "/id":
			newURL := "/"
			resp, _ := testRequestPost(t, ts, http.MethodGet, newURL, "")
			defer resp.Body.Close()
			assert.Equal(t, v.status, resp.StatusCode)
			assert.Equal(t, v.contentType, resp.Header.Get("Location"))
		case "/api/shorten":
			resp, _ := testRequestJSON(t, ts, v.method, v.api, v.url)
			defer resp.Body.Close()
			assert.Equal(t, v.status, resp.StatusCode)
			assert.Equal(t, v.contentType, resp.Header.Get("Content-Type"))
		}

	}

}

func TestGzipCompression(t *testing.T) {
	logger.Initialize()
	ts := httptest.NewServer(Router(&cfg, strg, strg.DB))
	defer ts.Close()

	fmt.Println("/")

	requestBody := `{
        "url": "https://yan.ru/"
    }`

	t.Run("sends_gzip", func(t *testing.T) {

		buf := bytes.NewBuffer(nil)
		zb := gzip.NewWriter(buf)
		_, err := zb.Write([]byte(requestBody))
		require.NoError(t, err)
		err = zb.Close()
		require.NoError(t, err)

		r := httptest.NewRequest("POST", ts.URL+"/api/shorten", buf)
		r.RequestURI = ""
		r.Header.Set("Content-Encoding", "gzip")

		resp, err := http.DefaultClient.Do(r)
		require.NoError(t, err)

		defer resp.Body.Close()

		_, err = io.ReadAll(resp.Body)
		require.NoError(t, err)

	})
}
