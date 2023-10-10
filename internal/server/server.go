package server

import (
	"log"
	"net/http"

	"github.com/poggerr/go_shortener/internal/logger"
)

// Server запускает сервер
func Server(addr string, hand http.Handler) {

	server := &http.Server{
		Addr:           addr,
		Handler:        hand,
		TLSConfig:      nil,
		MaxHeaderBytes: 16 * 1024,
	}

	logger.Initialize().Info("Running server on: ", addr)

	log.Fatal(server.ListenAndServe())
}
