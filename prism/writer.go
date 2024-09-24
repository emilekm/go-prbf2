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

func (w *Writer) WriteMessage(subject Subject, body []byte) error {
	ew := &errWriter{w: w.W}

	w.mutex.Lock()
	defer w.mutex.Unlock()

	ew.write(SeparatorStart)
	ew.write(stringToBytes(string(subject)))
	ew.write(SeparatorSubject)
	ew.write(body)
	ew.write(SeparatorEnd)
	ew.write(SeparatorNull)

	if ew.err != nil {
		return ew.err
	}

	return w.W.Flush()
}
