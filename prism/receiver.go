package prism

import (
	"bufio"
	"bytes"
	"io"
	"log/slog"
)

type Receiver struct {
	r      io.Reader
	broker *Broker[Message]
}

func NewReceiver(r io.Reader) *Receiver {
	receiver := &Receiver{
		r:      r,
		broker: NewBroker[Message](),
	}

	receiver.start()

	return receiver
}

func (r *Receiver) Subscribe() Subscriber[Message] {
	return r.broker.Subscribe()
}

func (r *Receiver) Unsubscribe(sub Subscriber[Message]) {
	r.broker.Unsubscribe(sub)
}

func (r *Receiver) start() {
	scanner := bufio.NewScanner(r.r)
	scanner.Split(splitMessages)

	go func() {
		for scanner.Scan() {
			buf := scanner.Bytes()
			msg, err := Decode(buf)
			if err != nil {
				slog.Warn("failed to decode message", "error", err)
			}
			r.broker.Publish(msg)
		}
	}()
}

func splitMessages(data []byte, atEOF bool) (advance int, token []byte, err error) {
	msg, _, found := bytes.Cut(data, SeparatorNull)
	if !found {
		return 0, nil, nil
	}

	advance = len(msg) + len(SeparatorNull)

	return advance, msg, nil
}
