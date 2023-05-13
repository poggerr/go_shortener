package handlers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/poggerr/go_shortener/internal/app/shorten"
	"github.com/poggerr/go_shortener/internal/config"
	"io"
	"log"
	"net/http"
	"os"
)

func GetPage(res http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	ans := shorten.UnShoring(id)
	res.Header().Set("Location", ans)
	res.WriteHeader(307)

}

func PostPage(res http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
		fmt.Println(err.Error())
		res.Write([]byte("Ошибка запроса"))
		return
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return
	}

	local := config.GetDefUrl()

	short, err := shorten.Shorting(string(body))

	if err != nil {
		http.Error(res, "Введите ссылку", http.StatusBadRequest)
		return
	}

	local += short

	res.Header().Set("content-type", "text/plain ")

	res.WriteHeader(http.StatusCreated)

	res.Write([]byte(local))

}
