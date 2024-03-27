package prism

import (
	"io"
)

type Transmitter struct {
	writer io.Writer
}

func NewTransmitter(w io.Writer) *Transmitter {
	return &Transmitter{
		writer: w,
	}
}

func (s *Transmitter) Send(subject Subject, fields ...[]byte) error {
	msg := NewMessage(subject, fields...)
	return s.SendRaw(msg.Encode())
}

func (s *Transmitter) SendRaw(msg []byte) error {
	_, err := s.writer.Write(msg)
	return err
}
