package prism

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
)

type Receiver struct {
	reader io.Reader
	msgCh  chan Message
}

func NewReceiver(r io.Reader) *Receiver {
	receiver := &Receiver{
		reader: r,
		msgCh:  make(chan Message),
	}

	receiver.Start()

	return receiver
}

func (r *Receiver) C() <-chan Message {
	return r.msgCh
}

func (r *Receiver) Start() {
	scanner := bufio.NewScanner(r.reader)
	scanner.Split(splitMessages)

	go func() {
		for scanner.Scan() {
			msg, err := DecodeMessage(scanner.Bytes())
			if err != nil {
				fmt.Println("decode error:", err)
				continue
			}

			select {
			case r.msgCh <- *msg:
			default:
			}
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
