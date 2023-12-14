package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	server := server.New(config.BuildPort(config.Config.Port), time.Millisecond*time.Duration(config.Config.Timeout))

	WOWstorage := storage.New(storage.WordsOfWisdom)

	challenge := app.NewChallenge(config.Config.Difficulty)

	requeststore := storage.NewRequestStore(storage.ShardKey)

	app := app.New(&server, WOWstorage, requeststore, challenge)

	err := app.Run(ctx)
	if err != nil {
		config.Logger.Fatalf("Error while starting service: %v", err)
	}
}
