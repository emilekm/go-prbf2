package prism

import (
	"bufio"
)

type Writer struct {
	W *bufio.Writer
}

func NewWriter(w *bufio.Writer) *Writer {
	return &Writer{W: w}
}

func (w *Writer) WriteMessage(msg Message) error {
	return w.Write(Encode(msg))
}

func (w *Writer) WriteRawMessage(rawMessage *RawMessage) error {
	return w.Write(rawMessage.Encode())
}
func (w *Writer) Write(data []byte) error {
	_, err := w.W.Write(data)
	if err != nil {
		return err
	}
	err = w.W.WriteByte(SeparatorNull1)
	if err != nil {
		return err
	}

	return w.W.Flush()
}
