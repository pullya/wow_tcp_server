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
	config.ReadConfig()
	config.InitLogger()
	log.SetLevel(config.Config.LogLevel.ToLogrusFormat())
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigCh
		config.Logger.Warnf("Received signal %v. Shutting down...", sig)

		cancel()
	}()

	server := server.New(config.Config.TcpPort)

	storage := storage.New(storage.WordsOfWisdom)

	challenge := app.NewChallenge(config.Config.Difficulty)

	app := app.New(&server, storage, challenge)

	err := app.Run(ctx)
	if err != nil {
		config.Logger.Fatalf("Error while starting service: %v", err)
	}
}
