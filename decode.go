package prdemo

import (
	"encoding/binary"
	"io"

	"github.com/ghostiam/binstruct"
)

type Decoder interface {
	Decode(DemoReader) error
}

type DemoReader interface {
	Decode(interface{}) error
}

type demoReader struct {
	decoder *binstruct.Decoder
}

func NewDemoReader(reader io.ReadSeeker) DemoReader {
	return &demoReader{
		decoder: binstruct.NewDecoder(reader, binary.LittleEndian),
	}
}

func (r *demoReader) Decode(v interface{}) error {
	if d, ok := v.(Decoder); ok {
		return d.Decode(r)
	}

	return r.decoder.Decode(v)
}
