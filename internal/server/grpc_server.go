package server

import (
	"github.com/poggerr/go_shortener/internal/service"
	"google.golang.org/grpc"
)

func NewGRPCServer(baseURL string, repo service.URLShortenerService) *grpc.Server {
	//linkStore := repo
	//handler := handlers.NewURLShortener(baseURL, linkStore)

	s := grpc.NewServer()
	// TODO добавить pb.RegisterShortenerServer
	//pb.RegisterShortenerServer(s, handler)
	return s
}
