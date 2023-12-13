package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/pullya/wow_tcp_server/tcp-client/internal/app"
	"github.com/pullya/wow_tcp_server/tcp-client/internal/client"
	"github.com/pullya/wow_tcp_server/tcp-client/internal/config"
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

	client := client.New(config.Config.Address)
	challenge := app.NewChallenge()

	app := app.New(&client, &challenge)

	app.Run(ctx)
}
