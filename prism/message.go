package prism

import (
	"bytes"
	"errors"
)

type Message struct {
	Subject Subject
	Content []byte
}

func NewMessage(subject Subject, fields ...[]byte) Message {
	return Message{
		Subject: subject,
		Content: bytes.Join(fields, SeparatorField),
	}
}

func (m *Message) Encode() []byte {
	return bytes.Join(
		[][]byte{
			SeparatorStart,
			[]byte(m.Subject),
			SeparatorSubject,
			m.Content,
			SeparatorEnd,
		},
		[]byte{},
	)
}

func (m *Message) Fields() [][]byte {
	return bytes.Split(m.Content, SeparatorField)
}

func DecodeMessage(data []byte) (*Message, error) {
	// Make sure we arent reading data after the end separator
	// or with it
	data, _, _ = bytes.Cut(data, SeparatorEnd)
	// Skipping data before the start separator
	// It might be empty spaces
	_, data, found := bytes.Cut(data, SeparatorStart)
	if !found {
		return nil, errors.New("missing start separator")
	}

	subject, content, found := bytes.Cut(data, SeparatorSubject)
	if !found {
		return nil, errors.New("missing subject separator")
	}

	return &Message{
		Subject: Subject(subject),
		Content: content,
	}, nil
}
