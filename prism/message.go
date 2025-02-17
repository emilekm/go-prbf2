package prism

import (
	"bytes"
	"errors"
)

// Message represents a message received from the server with header decoded.
type Message struct {
	subject Subject
	body    []byte
}

func NewMessage(subject Subject, body []byte) *Message {
	return &Message{
		subject: subject,
		body:    body,
	}
}

// Subject returns the subject of the message.
func (m *Message) Subject() Subject {
	return m.subject
}

// Body returns the body of the message.
// The whole slice is copied so it's safe to modify.
func (m *Message) Body() []byte {
	return m.body[:]
}

func (m *Message) MarshalBinary() ([]byte, error) {
	return bytes.Join([][]byte{
		SeparatorStart,
		[]byte(m.subject),
		SeparatorSubject,
		m.body,
		SeparatorEnd,
	}, []byte{}), nil
}

func (m *Message) UnmarshalBinary(data []byte) error {
	start := bytes.Index(data, SeparatorStart)
	subject := bytes.Index(data, SeparatorSubject)
	end := bytes.Index(data, SeparatorEnd)
	if start == -1 || end == -1 {
		return errors.New("prism: invalid message")
	}

	m.subject = Subject(data[start+1 : subject])
	m.body = data[subject+1 : end]

	return nil
}
