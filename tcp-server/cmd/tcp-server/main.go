package main

import (
	"context"
	"fmt"

	service "github.com/pullya/wow_tcp_server/tcp-server/internal/app"
	"github.com/pullya/wow_tcp_server/tcp-server/internal/config"
	"github.com/pullya/wow_tcp_server/tcp-server/internal/server"
	"github.com/pullya/wow_tcp_server/tcp-server/internal/storage"
)

func main() {
	ctx := context.Background()

	wowServer := server.NewTcpServer(config.TcpPort)

	wowStorage := storage.NewInMemStorage(storage.WordsOfWisdom)

	wowService := service.NewWowService(wowServer, wowStorage)

	err := wowService.Run(ctx)
	if err != nil {
		fmt.Println("Error while starting service:", err)
		return
	}
}
