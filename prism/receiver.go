package prism

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
)

type Receiver struct {
	reader io.Reader
	broker *Broker[Message]
}

func NewReceiver(r io.Reader) *Receiver {
	receiver := &Receiver{
		reader: r,
		broker: NewBroker[Message](),
	}

	receiver.Start()

	return receiver
}

func (r *Receiver) Listen() *Subscriber[Message] {
	return r.broker.Subscribe()
}

func (r *Receiver) Start() {
	scanner := bufio.NewScanner(io.TeeReader(r.reader, os.Stdout))
	scanner.Split(splitMessages)

	go func() {
		for scanner.Scan() {
			msg, err := DecodeMessage(scanner.Bytes())
			if err != nil {
				fmt.Println("decode error:", err)
				continue
			}

			r.broker.Publish(*msg)
		}
	}()
}

func splitMessages(data []byte, atEOF bool) (advance int, token []byte, err error) {
	msg, _, found := bytes.Cut(data, SeparatorEnd)
	if !found {
		return 0, nil, nil
	}

	advance = len(msg) + len(SeparatorEnd)

	return advance, msg, nil
}
