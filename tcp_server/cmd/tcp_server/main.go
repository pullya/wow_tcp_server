package main

import (
	"context"

	"github.com/pullya/wow_tcp_server/tcp_server/internal/server"
	"github.com/rs/zerolog/log"
)

func main() {
	ctx := context.Background()

	if err := server.RunServer(ctx); err != nil {
		log.Fatal().Str("service", "TCP-server").Msgf("Error while starting TCP-server: %v", err)
	}
}
