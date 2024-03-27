package prism

import (
	"errors"
	"fmt"
	"slices"
	"sync"
	"time"
)

type Responder struct {
	receiver     *Receiver
	chatReceiver *BufferReceiver[ChatMessage]
	sender       *Transmitter
	mutex        sync.Mutex
}

func NewResponder(receiver *Receiver, chatReceiver *BufferReceiver[ChatMessage], sender *Transmitter) *Responder {
	return &Responder{
		receiver:     receiver,
		chatReceiver: chatReceiver,
		sender:       sender,
	}
}

type SendOpts struct {
	ResponseSubjects []Subject
	ChatMessageType  *ChatMessageType
}

type Response struct {
	Messages    []Message
	ChatMessage *ChatMessage
}

func (r *Responder) Send(msg Message, opts *SendOpts) (*Response, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if opts == nil {
		opts = &SendOpts{}
	}

	subjects := append([]Subject{}, opts.ResponseSubjects...)

	var resp Response
	var subjectErr error

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		timer := time.NewTimer(5 * time.Second)
		for {
			if len(subjects) == 0 {
				return
			}

			sub := r.receiver.Listen()

			select {
			case <-timer.C:
				subjectErr = errors.New("timeout")
				return
			case m := <-sub.Channel:
				println(m.Subject)
				if i := slices.Index(subjects, m.Subject); i != -1 {
					resp.Messages = append(resp.Messages, m)
					subjects = append(subjects[:i], subjects[i+1:]...)
				}
				if slices.Contains(errorSubjects, m.Subject) {
					var msgErr2 Error
					err := UnmarshalInto(msg, &msgErr2)
					if err != nil {
						subjectErr = fmt.Errorf("unmarshal error: %w", err)
						return
					}
					subjectErr = msgErr2
					return
				}
			}
		}
	}()

	var chatErr error

	go func() {
		defer wg.Done()

		timer := time.NewTimer(2 * time.Second)
		sub := r.chatReceiver.Listen()

		for {
			if opts.ChatMessageType == nil {
				return
			}

			select {
			case <-timer.C:
				chatErr = errors.New("timeout")
				return
			case m := <-sub.Channel:
				if m.Type == *opts.ChatMessageType {
					resp.ChatMessage = &m
					return
				}
			}
		}
	}()

	time.Sleep(300 * time.Millisecond)

	err := r.sender.SendRaw(msg.Encode())
	if err != nil {
		return nil, fmt.Errorf("send: %w", err)
	}

	wg.Wait()

	if subjectErr != nil {
		return nil, fmt.Errorf("send subject err: %w", subjectErr)
	}

	if chatErr != nil {
		return nil, fmt.Errorf("send chat err: %w", chatErr)
	}

	return &resp, nil
}
