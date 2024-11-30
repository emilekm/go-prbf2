/*
PRISM is a PR:BF2 server moderation tool.
The server sends player updates to all connected clients, but also allows for running queries (commands).

This client aims to provide a robust interface for interacting with the PRISM server,
whether it be for reading player updates or sending commands to the server.
*/
package prism2

import (
	"bufio"
	"io"
	"net"
	"net/textproto"
)

type Client struct {
	Reader
	Writer

	pipeline textproto.Pipeline
	conn     io.ReadWriteCloser
}

func NewClient(conn io.ReadWriteCloser) *Client {
	return &Client{
		Reader:   Reader{R: bufio.NewReader(conn)},
		Writer:   Writer{W: bufio.NewWriter(conn)},
		pipeline: textproto.Pipeline{},
		conn:     conn,
	}
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

// Next returns the next id for a request.
func (c *Client) Next() uint {
	return c.pipeline.Next()
}

// StartRequest blocks until it is time to send the request with the given id.
func (c *Client) StartRequest(id uint) {
	c.pipeline.StartRequest(id)
}

// EndRequest notifies p that the request with the given id has been sent.
func (c *Client) EndRequest(id uint) {
	c.pipeline.EndRequest(id)
}
