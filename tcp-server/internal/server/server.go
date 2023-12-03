package server

import (
	"context"
	"fmt"
	"net"
)

//go:generate mockery --name=IServer --output=mocks --case=underscore
type IServer interface {
	RunServer(ctx context.Context) (net.Listener, error)
}

type TcpServer struct {
	Port string
}

func NewTcpServer(port string) TcpServer {
	return TcpServer{
		Port: port,
	}
}

func (ts TcpServer) RunServer(ctx context.Context) (net.Listener, error) {
	fmt.Println("Launching tcp-server...")

	listener, err := net.Listen("tcp", ts.Port)
	if err != nil {
		return nil, err
	}

	return listener, nil
}
