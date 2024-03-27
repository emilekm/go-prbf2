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
	config ClientConfig
	conn   net.Conn

	receiver      *Receiver
	chatReceiver  *BufferReceiver[ChatMessage]
	killsReceiver *BufferReceiver[KillMessage]
	sender        *Transmitter

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

	c.receiver = NewReceiver(c.conn)
	c.chatReceiver = NewBufferReceiver[ChatMessage](c.receiver)
	c.killsReceiver = NewBufferReceiver[KillMessage](c.receiver)
	c.sender = NewTransmitter(c.conn)

	c.Responder = NewResponder(c.receiver, c.chatReceiver, c.sender)

	return auth(c.Responder, c.config)
}

func (c *Client) Disconnect() {
	c.conn.Close()
}
