package prism

import (
	"net"
	"time"
)

const (
	readTimeout  = 5 * time.Minute
	writeTimeout = 5 * time.Minute
)

type Client struct {
	config    ClientConfig
	conn      net.Conn
	receiver  *MessageReceiver
	sender    *MessageSender
	Responder *Responder
}

type ClientConfig struct {
	IP   string
	Port string
	User string
	Pass string
}

func NewClient(config ClientConfig) *Client {
	return &Client{
		config: config,
	}
}

func (c *Client) Connect() error {
	conn, err := net.Dial("tcp", net.JoinHostPort(c.config.IP, c.config.Port))
	if err != nil {
		return err
	}

	c.conn = conn
	c.conn.SetDeadline(time.Time{})

	c.receiver = NewMessageReceiver(c.conn)
	c.sender = NewMessageSender(c.conn)

	c.Responder = NewResponder(c.receiver, c.sender)

	return auth(c.Responder, c.config)
}

func (c *Client) Disconnect() {
	c.conn.Close()
}
