package prism

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
)

type MessageReceiver struct {
	reader io.Reader
	msgCh  chan Message
}

func NewMessageReceiver(r io.Reader) *MessageReceiver {
	receiver := &MessageReceiver{
		reader: r,
		msgCh:  make(chan Message),
	}

	receiver.Start()

	return receiver
}

func (mr *MessageReceiver) C() <-chan Message {
	return mr.msgCh
}

func (mr *MessageReceiver) Start() {
	scanner := bufio.NewScanner(mr.reader)
	scanner.Split(splitMessages)

	go func() {
		for scanner.Scan() {
			msg, err := DecodeMessage(scanner.Bytes())
			if err != nil {
				fmt.Println("decode error:", err)
				continue
			}

			select {
			case mr.msgCh <- *msg:
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
