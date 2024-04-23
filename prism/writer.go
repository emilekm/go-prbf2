package prism

import (
	"bufio"
	"io"
	"sync"
)

type Writer struct {
	W     *bufio.Writer
	mutex sync.Mutex
}

func NewWriter(w *bufio.Writer) *Writer {
	return &Writer{W: w}
}

func (w *Writer) WriteMessage(msg Message) error {
	rawMsg, err := EncodeMessage(msg)
	if err != nil {
		return err
	}

	return w.WriteRawMessage(rawMsg)
}

type errWriter struct {
	w   io.Writer
	err error
}

func (ew *errWriter) write(buf []byte) {
	if ew.err != nil {
		return
	}
	_, ew.err = ew.w.Write(buf)
}

func (w *Writer) WriteRawMessage(rawMessage *RawMessage) error {
	ew := &errWriter{w: w.W}

	w.mutex.Lock()
	defer w.mutex.Unlock()

	ew.write(SeparatorStart)
	ew.write(stringToBytes(string(rawMessage.Subject())))
	ew.write(SeparatorSubject)
	ew.write(rawMessage.Content())
	ew.write(SeparatorEnd)
	ew.write(SeparatorNull)

	if ew.err != nil {
		return ew.err
	}

	return w.W.Flush()
}
