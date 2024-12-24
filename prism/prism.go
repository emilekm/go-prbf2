/*
PRISM is a PR:BF2 server moderation tool.
The server sends player updates to all connected clients, but also allows for running queries (commands).

This client aims to provide a robust interface for interacting with the PRISM server,
whether it be for reading player updates or sending commands to the server.
*/
package prism

import (
	"bufio"
	"io"
	"net"
	"net/textproto"
)

type Client struct {
	Reader
	Writer

	*broker

	Server   *serverService
	Gameplay *gameplayService
	Players  *playersService
	Users    *usersService
	Admin    *adminService

	textproto.Pipeline
	conn io.ReadWriteCloser
}

func NewClient(conn io.ReadWriteCloser) *Client {
	c := &Client{
		Reader:   Reader{R: bufio.NewReader(conn)},
		Writer:   Writer{W: bufio.NewWriter(conn)},
		Pipeline: textproto.Pipeline{},
		conn:     conn,
	}

	c.Server = &serverService{c: c}
	c.Gameplay = &gameplayService{c: c}
	c.Players = &playersService{c: c}
	c.Users = &usersService{c: c}
	c.Admin = &adminService{c: c}

	c.broker = newBroker(c)

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
	c.broker.Close()

	return c.conn.Close()
}
