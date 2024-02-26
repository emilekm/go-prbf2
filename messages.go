package prdemo

import (
	"io"

	"github.com/ghostiam/binstruct"
)

type Message struct {
	Type MessageType
	r    binstruct.Reader
}

func NewMessage(r io.ReadSeeker) (*Message, error) {
	br := newBinReader(r)
	typ, err := br.ReadUint8()
	if err != nil {
		return nil, err
	}

	return &Message{
		Type: MessageType(typ),
		r:    br,
	}, nil
}

func (m *Message) Decode(v interface{}) error {
	if d, ok := v.(Decoder); ok {
		return d.Decode(m)
	}

	return m.walk(v)
}
