package main

import (
	"fmt"
	"github.com/poggerr/go_shortener/internal/app"
	"io"
	"net/http"
	"strings"
)

func main() {
	mux := http.NewServeMux()
	//mux.HandleFunc(`/`, postPage)
	//router.HandleFunc(`/{id}`, getPage)
	mux.HandleFunc("/", mainHendler)
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Println(err)
	}
}

func mainHendler(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		postPage(res, req)
	case "GET":
		id := req.URL.Path
		getPage(res, id)
	}

}

func getPage(res http.ResponseWriter, id string) {
	//if req.Method != http.MethodGet {
	//	http.Error(res, "Only GET requests are allowed!", http.StatusBadRequest)
	//	return
	//}
	s := strings.Split(id, "/")
	ans := app.UnShorting(s[1])

	res.Header().Set("Location", ans)
	res.WriteHeader(307)

}

func postPage(res http.ResponseWriter, req *http.Request) {
	//if req.Method != http.MethodPost {
	//	http.Error(res, "Only Post requests are allowed!", http.StatusBadRequest)
	//	return
	//}

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
