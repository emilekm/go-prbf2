package prism

import (
	"bufio"
	"sync"
)

type Writer struct {
	W *bufio.Writer

	mutex sync.Mutex
}

func NewWriter(w *bufio.Writer) *Writer {
	return &Writer{W: w}
}

func (w *Writer) WriteMessage(msg *Message) error {
	content, err := msg.MarshalBinary()
	if err != nil {
		return err
	}

	return w.WriteMessageBytes(content)
}

func (w *Writer) WriteMessageBytes(content []byte) error {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	_, err := w.W.Write(content)
	if err != nil {
		return err
	}

	_, err = w.W.Write(SeparatorNull)
	if err != nil {
		return err
	}

	return w.W.Flush()
}
