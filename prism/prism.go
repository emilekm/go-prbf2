package prism

import (
	"bufio"
	"io"
	"net"
	"net/textproto"
)

type Client struct {
	Receiver
	Responder
	Auth *Auth

	conn     io.ReadWriteCloser
	reader   *Reader
	writer   *Writer
	pipeline *textproto.Pipeline
}

func NewClient(conn io.ReadWriteCloser) *Client {
	pipeline := &textproto.Pipeline{}
	reader := &Reader{R: bufio.NewReader(conn)}
	writer := &Writer{W: bufio.NewWriter(conn)}

	receiver := NewReceiver(reader)

	c := &Client{
		Receiver:  *receiver,
		Responder: *NewResponder(receiver, writer, pipeline),

		conn:     conn,
		reader:   reader,
		writer:   writer,
		pipeline: pipeline,
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
