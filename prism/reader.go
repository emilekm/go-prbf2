package prism

import (
	"bufio"
)

type Reader struct {
	r *bufio.Reader
}

func NewReader(r *bufio.Reader) *Reader {
	return &Reader{r: r}
}

func (r *Reader) ReadMessage() (*Message, error) {
	buf, err := r.ReadMessageBytes()
	if err != nil {
		return nil, err
	}

	var msg Message
	err = msg.UnmarshalBinary(buf)
	if err != nil {
		return nil, err
	}

	return &msg, nil
}

func (r *Reader) ReadMessageBytes() ([]byte, error) {
	return r.r.ReadBytes(SeparatorNull1)
}
