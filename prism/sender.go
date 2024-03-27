package prism

import "io"

type MessageSender struct {
	writer io.Writer
}

func NewMessageSender(w io.Writer) *MessageSender {
	return &MessageSender{
		writer: w,
	}
}

func (ms *MessageSender) Send(subject Subject, fields ...[]byte) error {
	msg := NewMessage(subject, fields...)
	return ms.SendRaw(msg.Encode())
}

func (ms *MessageSender) SendRaw(msg []byte) error {
	_, err := ms.writer.Write(msg)
	return err
}
