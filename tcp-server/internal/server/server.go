package server

import (
	"bufio"
	"context"
	"net"

	"github.com/pullya/wow_tcp_server/tcp-server/internal/config"
)

//go:generate mockery --name=ServerProvider --output=mocks --case=underscore
type ServerProvider interface {
	Run(ctx context.Context) (net.Listener, error)
	SendMessage(ctx context.Context, conn net.Conn, mess []byte) error
	ReceiveMessage(ctx context.Context, conn net.Conn) (string, error)
}

type Server struct {
	Port string
	Conn net.Conn
}

func New(port string) Server {
	return Server{
		Port: port,
	}
}

func (ts *Server) Run(ctx context.Context) (net.Listener, error) {
	config.Logger.Info("Launching tcp-server...")

	listener, err := net.Listen("tcp", ts.Port)
	if err != nil {
		return nil, err
	}

	return listener, nil
}

func (ts *Server) SendMessage(ctx context.Context, conn net.Conn, mess []byte) error {
	if _, err := conn.Write(mess); err != nil {
		return err
	}

	return nil
}

func (ts *Server) ReceiveMessage(ctx context.Context, conn net.Conn) (string, error) {
	return bufio.NewReader(conn).ReadString('\n')
}
