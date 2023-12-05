package server

import (
	"bufio"
	"context"
	"net"

	"github.com/pullya/wow_tcp_server/tcp-server/internal/config"
	log "github.com/sirupsen/logrus"
)

type IServer interface {
	RunServer(ctx context.Context) (net.Listener, error)
	SendMessage(ctx context.Context, mess []byte) error
	ReceiveMessage(ctx context.Context) (string, error)
	SetConn(conn net.Conn)
	CloseConn()
}

type TcpServer struct {
	Port string
	Conn net.Conn
}

func NewTcpServer(port string) TcpServer {
	return TcpServer{
		Port: port,
	}
}

func (ts *TcpServer) RunServer(ctx context.Context) (net.Listener, error) {
	log.WithField("service", config.ServiceName).Info("Launching tcp-server...")

	listener, err := net.Listen("tcp", ts.Port)
	if err != nil {
		return nil, err
	}

	return listener, nil
}

func (ts *TcpServer) SendMessage(ctx context.Context, mess []byte) error {
	if _, err := ts.Conn.Write(mess); err != nil {
		return err
	}

	return nil
}

func (ts *TcpServer) ReceiveMessage(ctx context.Context) (string, error) {
	return bufio.NewReader(ts.Conn).ReadString('\n')
}

func (ts *TcpServer) SetConn(conn net.Conn) {
	ts.Conn = conn
}

func (ts *TcpServer) CloseConn() {
	connection := ts.Conn
	connection.Close()
}
