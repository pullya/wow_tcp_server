package service

import (
	"bufio"
	"context"
	"fmt"
	"net"

	"github.com/pullya/wow_tcp_server/tcp-server/internal/server"
	"github.com/pullya/wow_tcp_server/tcp-server/internal/storage"
)

var connCnt = 0

type WowService struct {
	Server  server.IServer
	Storage storage.IStorage
}

func NewWowService(tcpServer server.IServer, storage storage.IStorage) WowService {
	return WowService{
		Server:  tcpServer,
		Storage: storage,
	}
}

func (ws WowService) Run(ctx context.Context) error {
	listener, err := ws.Server.RunServer(ctx)
	if err != nil {
		fmt.Println("Error while starting TCP-server:", err)
	}

	defer listener.Close()

	fmt.Println("Waiting for connection...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		connCnt++
		fmt.Printf("[conn_%d]New connection established!\n", connCnt)

		go ws.HandleConnection(ctx, conn, connCnt)
	}
}

func (ws WowService) HandleConnection(ctx context.Context, conn net.Conn, id int) {
	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Printf("[conn_%d]Error while reading incoming request: %v\n", id, err)
			return
		}

		fmt.Printf("[conn_%d]Message received: %s\n", id, string(message))

		newmessage := ws.Storage.GetRandomWoW(ctx)

		fmt.Printf("[conn_%d]Prepared response: %s\n", id, string(newmessage))

		_, err = conn.Write([]byte(newmessage + "\n"))
		if err != nil {
			fmt.Printf("[conn_%d]Error while sending response: %v\n", id, err)
			return
		}
	}
}
