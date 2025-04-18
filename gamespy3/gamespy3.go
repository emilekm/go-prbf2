package gamespy3

import (
	"net"
	"sync"
)

type Client struct {
	mutex sync.Mutex
	conn  *net.UDPConn
}

func New(conn *net.UDPConn) *Client {
	return &Client{
		conn: conn,
	}
}
