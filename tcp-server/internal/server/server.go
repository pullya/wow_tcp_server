package server

import (
	"context"
	"net"

	"github.com/rs/zerolog/log"
)

type Server interface {
	RunServer(ctx context.Context) (net.Conn, error)
}

type TcpServer struct {
	Port string
}

func NewTcpServer(port string) TcpServer {
	return TcpServer{
		Port: port,
	}
}

func (ts TcpServer) RunServer(ctx context.Context) (net.Conn, error) {
	log.Info().Msg("Launching tcp-server...")

	listener, err := net.Listen("tcp", ts.Port)
	if err != nil {
		return nil, err
	}

	conn, err := listener.Accept()
	if err != nil {
		return nil, err
	}

	return conn, nil
}
