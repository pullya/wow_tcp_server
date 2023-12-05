package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/pullya/wow_tcp_server/tcp-server/internal/app"
	"github.com/pullya/wow_tcp_server/tcp-server/internal/config"
	"github.com/pullya/wow_tcp_server/tcp-server/internal/server"
	"github.com/pullya/wow_tcp_server/tcp-server/internal/storage"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(config.LogLevel)
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigCh
		log.WithField("service", config.ServiceName).Warnf("Received signal %v. Shutting down...", sig)

		cancel()
	}()

	wowServer := server.NewTcpServer(config.TcpPort)

	wowStorage := storage.NewInMemStorage(storage.WordsOfWisdom)

	wowChallenge := app.NewChallenge(config.PowDifficulty)

	wowService := app.NewWowService(&wowServer, wowStorage, wowChallenge)

	err := wowService.Run(ctx)
	if err != nil {
		log.WithField("service", config.ServiceName).Fatalf("Error while starting service: %v", err)
	}
}
