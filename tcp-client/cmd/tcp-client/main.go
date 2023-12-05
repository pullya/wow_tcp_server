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

	wowClient := client.NewClient(config.Address)
	wowChallenge := app.NewChallenge()

	wowService := app.NewWowService(&wowClient, &wowChallenge)

	wowService.Run(ctx)
}
