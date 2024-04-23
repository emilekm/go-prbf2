package prism

import (
	"bufio"
	"bytes"
	"errors"
)

type Reader struct {
	R *bufio.Reader
}

func NewReader(r *bufio.Reader) *Reader {
	return &Reader{R: r}
}

func decodeData(data []byte) (*RawMessage, error) {
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

	return NewRawMessage(Subject(subject), content), nil
}

func (r *Reader) ReadMessage() (*RawMessage, error) {
	buf, err := r.readMessageBytes()
	if err != nil {
		return nil, err
	}

	return decodeData(buf)
}

func (r *Reader) readMessageBytes() ([]byte, error) {
	return r.R.ReadBytes(SeparatorNull1)
}
