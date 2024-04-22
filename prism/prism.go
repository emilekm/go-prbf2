package prism

import (
	"bufio"
	"io"
	"net"
	"net/textproto"
)

type Client struct {
	*textproto.Pipeline
	Reader
	Writer
	Auth *Auth
	conn io.ReadWriteCloser
}

func NewClient(conn io.ReadWriteCloser) *Client {
	pipeline := &textproto.Pipeline{}

	c := &Client{
		Pipeline: pipeline,
		Reader:   Reader{R: bufio.NewReader(conn)},
		Writer:   Writer{W: bufio.NewWriter(conn)},
		conn:     conn,
	}

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
