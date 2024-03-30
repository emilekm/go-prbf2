package prism

import (
	"io"
	"net/textproto"
	"time"
)

type Transmitter struct {
	pipeline *textproto.Pipeline
	writer   io.Writer
}

func NewTransmitter(w io.Writer, pipeline *textproto.Pipeline) *Transmitter {
	return &Transmitter{
		pipeline: pipeline,
		writer:   w,
	}
}

func (s *Transmitter) Send(msg Message) error {
	buf, err := Encode(msg)
	if err != nil {
		return err
	}

	return s.SendRaw(buf)
}

func (s *Transmitter) SendRaw(msg []byte) error {
	id := s.pipeline.Next()
	s.pipeline.StartRequest(id)
	_, err := s.writer.Write(msg)
	// Sleep for 1 seoncd to allow any responses buffer to be sent
	// so Responder.SendWithRespone won't get invalid response
	time.Sleep(1 * time.Second)
	s.pipeline.EndRequest(id)
	return err
}
