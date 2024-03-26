package prism

import (
	"bytes"
	"errors"
)

type Message struct {
	Subject Subject
	Fields  [][]byte
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
	data, found := bytes.CutPrefix(data, SeparatorStart)
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
