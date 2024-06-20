package prism

import (
	"bufio"
	"context"
	"io"
	"net"
	"net/textproto"
)

type Message interface {
	Subject() Subject
}

type RawMessage struct {
	subject Subject
	body    []byte
}

func NewRawMessage(subject Subject, content []byte) *RawMessage {
	return &RawMessage{
		subject: subject,
		body:    content,
	}
}

func (m RawMessage) Subject() Subject {
	return m.subject
}

func (m RawMessage) Body() []byte {
	return m.body
}

func (m RawMessage) MarshalMessage() ([]byte, error) {
	return m.body, nil
}

func (m *RawMessage) UnmarshalMessage(content []byte) error {
	m.body = content[:]
	return nil
}

type Client struct {
	textproto.Pipeline
	Reader
	Writer

	conn io.ReadWriteCloser
}

func NewClient(conn io.ReadWriteCloser) *Client {
	return &Client{
		Pipeline: textproto.Pipeline{},
		Reader:   Reader{R: bufio.NewReader(conn)},
		Writer:   Writer{W: bufio.NewWriter(conn)},
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

func (c *Client) Command(ctx context.Context, cmd Command, payload any, success Subject) (*RawMessage, error) {
	content := make([]byte, 0)
	if payload != nil {
		var err error
		content, err = MarshalMessage(payload)
		if err != nil {
			return nil, err
		}
	}

	id := c.Next()
	c.StartRequest(id)
	c.StartResponse(id)

	err := c.WriteMessage(NewRawMessage(cmd, content))
	c.EndRequest(id)
	if err != nil {
		c.EndResponse(id)
		return nil, err
	}

	msg, err := c.waitForMessage(ctx, success)
	c.EndResponse(id)
	if err != nil {
		return nil, err
	}

	return msg, nil
}
