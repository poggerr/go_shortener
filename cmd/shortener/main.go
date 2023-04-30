package main

import (
	"github.com/gorilla/mux"
	"github.com/poggerr/go_shortener/internal/app"
	"io"
	"net/http"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	router := mux.NewRouter()
	router.HandleFunc(`/`, postPage)
	router.HandleFunc(`/{id}`, getPage)
	return http.ListenAndServe(`:8080`, router)
}

func getPage(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]
	if req.Method != http.MethodGet {
		http.Error(res, "Only GET requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	ans := app.UnShorting(id)

	res.Header().Set("content-type", "text/plain")

	res.WriteHeader(http.StatusTemporaryRedirect)

	res.Header().Set("Location", "https://practicum.yandex.ru/ ")

	res.Write([]byte(ans))
}

func postPage(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Only Post requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	if err := req.ParseForm(); err != nil {
		res.Write([]byte(err.Error()))
		return
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return
	}

	local := "http://localhost:8080/"

	local += app.Shorting(string(body))

	res.Header().Set("content-type", "text/plain ")

	res.WriteHeader(http.StatusCreated)

	res.Write([]byte(local))

}
