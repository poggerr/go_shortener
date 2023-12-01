package server

import (
	"github.com/poggerr/go_shortener/internal/handlers"
	"github.com/poggerr/go_shortener/internal/service"
	"google.golang.org/grpc"

	pb "github.com/poggerr/go_shortener/internal/proto"
)

func NewGRPCServer(baseURL string, repo service.URLShortenerService) *grpc.Server {
	linkStore := repo
	handler := handlers.NewURLShortener(baseURL, linkStore)

	s := grpc.NewServer()
	pb.RegisterShortenerServer(s, handler)
	return s
}
