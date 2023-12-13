package server

import (
	"bufio"
	"context"
	"net"
	"time"

	"github.com/pullya/wow_tcp_server/tcp-server/internal/config"
)

//go:generate mockery --name=ServerProvider --output=mocks --case=underscore
type ServerProvider interface {
	Run(ctx context.Context) (net.Listener, error)
	SendMessage(ctx context.Context, conn net.Conn, mess []byte) error
	ReceiveMessage(ctx context.Context, conn net.Conn) (string, error)
	GetTimeout() time.Duration
}

type Server struct {
	Port    string
	Timeout time.Duration
}

func New(port string, timeout time.Duration) Server {
	return Server{
		Port:    port,
		Timeout: timeout,
	}
}

func (ts *Server) Run(_ context.Context) (net.Listener, error) {
	config.Logger.Info("Launching tcp-server...")

	listener, err := net.Listen("tcp", ts.Port)
	if err != nil {
		return nil, err
	}

	return listener, nil
}

func (ts *Server) SendMessage(_ context.Context, conn net.Conn, mess []byte) error {
	if _, err := conn.Write(mess); err != nil {
		return err
	}

	return nil
}

func (ts *Server) ReceiveMessage(_ context.Context, conn net.Conn) (string, error) {
	return bufio.NewReader(conn).ReadString('\n')
}

func (ts *Server) GetTimeout() time.Duration {
	return ts.Timeout
}
