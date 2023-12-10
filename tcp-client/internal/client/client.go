package client

import (
	"bufio"
	"context"
	"net"
)

//go:generate mockery --name=IClient --output=mocks --case=underscore
type IClient interface {
	RunClient(ctx context.Context) (net.Conn, error)
	SendMessage(ctx context.Context, conn net.Conn, mess []byte) error
	ReceiveMessage(ctx context.Context, conn net.Conn) (string, error)
	CloseConn(conn net.Conn)
}

type Client struct {
	Address string
}

func NewClient(addr string) Client {
	return Client{
		Address: addr,
	}
}

func (c *Client) RunClient(ctx context.Context) (net.Conn, error) {
	conn, err := net.Dial("tcp", c.Address)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func (c *Client) SendMessage(ctx context.Context, conn net.Conn, mess []byte) error {
	if _, err := conn.Write(mess); err != nil {
		return err
	}

	return nil
}

func (c *Client) ReceiveMessage(ctx context.Context, conn net.Conn) (string, error) {
	return bufio.NewReader(conn).ReadString('\n')
}

func (c *Client) CloseConn(conn net.Conn) {
	conn.Close()
}
