package main

import (
	"context"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/poggerr/go_shortener/internal/config"
	"github.com/poggerr/go_shortener/internal/server"
	"github.com/poggerr/go_shortener/internal/service"
	"github.com/poggerr/go_shortener/internal/utils"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
	cfg          *config.Config
	repo         service.URLShortenerService
)

func main() {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)
	srv := CreateServer()
	g := CreateGRPCServer()
	Run(srv, g)
}

func CreateServer() *http.Server {
	return server.Server(cfg.Serv, cfg.DefURL, repo, cfg.TrustedSubnet)
}

func CreateGRPCServer() *grpc.Server {
	return server.NewGRPCServer(cfg.DefURL, repo)
}

func Run(srv *http.Server, grpc *grpc.Server) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer stop()

	const (
		cert = "cert.pem"
		key  = "key.pem"
	)

	g, gCtx := errgroup.WithContext(ctx)

	log.Info().Msg("Start GRPC server...")

	g.Go(func() error {
		listen, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatal().Msgf("listen: %+v\n", err)
		}
		if err = grpc.Serve(listen); err != nil {
			log.Fatal().Msgf("serve: %+v\n", err)
		}
		return err
	})

	g.Go(func() error {
		if cfg.EnableHTTPS {
			log.Info().Msg("HTTPS enabled")
			err := utils.CreateTLSCert(cert, key)
			if err != nil {
				log.Fatal().Msgf("cert creation: %+v\n", err)
			}
			return srv.ListenAndServeTLS(cert, key)

		}
		log.Info().Msg("HTTPS is not enabled")
		return srv.ListenAndServe()
	})
	g.Go(func() error {
		<-gCtx.Done()
		grpc.Stop()
		return srv.Shutdown(context.Background())
	})

	if err := g.Wait(); err != nil {
		repo.Close()
		fmt.Printf("exit reason: %s \n", err)
	}
}
