package prism

import (
	"bytes"
	"errors"
)

type Message struct {
	Subject Subject
	Fields  [][]byte
}

func NewMessage(subject Subject, fields ...[]byte) Message {
	return Message{
		Subject: subject,
		Fields:  fields,
	}
}

func (m *Message) Encode() []byte {
	return bytes.Join(
		[][]byte{
			SeparatorStart,
			[]byte(m.Subject),
			SeparatorSubject,
			bytes.Join(m.Fields, SeparatorField),
			SeparatorEnd,
		},
		[]byte{},
	)
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

	subject, fields, found := bytes.Cut(data, SeparatorSubject)
	if !found {
		return nil, errors.New("missing subject separator")
	}

	return &Message{
		Subject: Subject(subject),
		Fields:  bytes.Split(fields, SeparatorField),
	}, nil
}
