package client

import (
	"bufio"
	"context"
	"net"
)

//go:generate mockery --name=IClient --output=mocks --case=underscore
type IClient interface {
	RunClient(ctx context.Context) error
	SendMessage(ctx context.Context, mess []byte) error
	ReceiveMessage(ctx context.Context) (string, error)
	CloseConn()
}

type Client struct {
	Address string
	Conn    net.Conn
}

func NewClient(addr string) Client {
	return Client{
		Address: addr,
	}
}

func (c *Client) RunClient(ctx context.Context) error {
	conn, err := net.Dial("tcp", c.Address)
	if err != nil {
		return err
	}
	c.Conn = conn
	return nil
}

func (c *Client) SendMessage(ctx context.Context, mess []byte) error {
	if _, err := c.Conn.Write(mess); err != nil {
		return err
	}

	return nil
}

func (c *Client) ReceiveMessage(ctx context.Context) (string, error) {
	return bufio.NewReader(c.Conn).ReadString('\n')
}

func (c *Client) CloseConn() {
	c.Conn.Close()
}
