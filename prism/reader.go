package prism

import (
	"bufio"
)

type Reader struct {
	R *bufio.Reader
}

func NewReader(r *bufio.Reader) *Reader {
	return &Reader{R: r}
}

func (r *Reader) ReadMessage() (*RawMessage, error) {
	buf, err := r.ReadMessageBytes()
	if err != nil {
		return nil, err
	}

	return DecodeRaw(buf)
}

func (r *Reader) ReadMessageBytes() ([]byte, error) {
	return r.R.ReadBytes(SeparatorNull1)
}
