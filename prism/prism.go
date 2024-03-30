package prism

import (
	"io"
	"net"
	"net/textproto"
)

type Client struct {
	*textproto.Pipeline
	Receiver
	Transmitter
	Responder
	Auth *Auth
	conn io.ReadWriteCloser
}

func NewClient(conn io.ReadWriteCloser) *Client {
	pipeline := &textproto.Pipeline{}

	c := &Client{
		Pipeline:    pipeline,
		Receiver:    *NewReceiver(conn),
		Transmitter: *NewTransmitter(conn, pipeline),
		conn:        conn,
	}

	c.Responder = *NewResponder(c)
	c.Auth = NewAuth(c)

	return c
}

func Dial(addrS string) (*Client, error) {
	addr, err := net.ResolveTCPAddr("tcp", addrS)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return nil, err
	}

	return NewClient(conn), nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}
