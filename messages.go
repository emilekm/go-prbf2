package prdemo

import (
	"encoding/binary"
	"io"

	"github.com/ghostiam/binstruct"
)

type Message struct {
	Type MessageType
	r    binstruct.Reader
}

func NewMessage(r io.ReadSeeker) (*Message, error) {
	var typ MessageType

	err := binary.Read(r, demoEndian, &typ)
	if err != nil {
		return nil, err
	}

	return &Message{
		Type: typ,
		r:    newBinReader(r),
	}, nil
}

func (m *Message) Decode(v interface{}) error {
	if d, ok := v.(Decoder); ok {
		return d.Decode(m)
	}

	unmarshaler := unmarshal{m.r}

	return unmarshaler.Unmarshal(v)
}
