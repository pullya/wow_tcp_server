package service

import (
	"bufio"
	"context"

	"github.com/pullya/wow_tcp_server/tcp-server/internal/server"
	"github.com/pullya/wow_tcp_server/tcp-server/internal/storage"
	"github.com/rs/zerolog/log"
)

type WowService struct {
	Server  server.Server
	Storage storage.Storage
}

func NewWowService(tcpServer server.Server, storage storage.Storage) WowService {
	return WowService{
		Server:  tcpServer,
		Storage: storage,
	}
}

func (ws WowService) Run(ctx context.Context) error {
	conn, err := ws.Server.RunServer(ctx)
	if err != nil {
		log.Error().Str("service", "TCP-server").Msgf("Error while starting TCP-server: %v", err)
	}

	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			return err
		}

		log.Info().Msgf("Message received: %s", string(message))

		newmessage := ws.Storage.GetRandomWoW(ctx)
		log.Info().Msgf("Prepared response: %s", string(newmessage))

		_, err = conn.Write([]byte(newmessage + "\n"))
		if err != nil {
			log.Error().Msgf("Error while sending response: %v", err)
			return err
		}
	}
}
