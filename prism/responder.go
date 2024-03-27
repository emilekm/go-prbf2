package prism

import (
	"errors"
	"fmt"
	"slices"
	"sync"
	"time"
)

type Responder struct {
	receiver *Receiver
	sender   *Transmitter
	mutex    sync.Mutex
}

func NewResponder(receiver *Receiver, sender *Transmitter) *Responder {
	return &Responder{
		receiver: receiver,
		sender:   sender,
	}
}

type SendOpts struct {
	ResponseSubjects []Subject
}

type Response struct {
	Messages []Message
}

func (r *Responder) Send(msg Message, opts *SendOpts) (*Response, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if opts == nil {
		opts = &SendOpts{}
	}

	subjects := append([]Subject{}, opts.ResponseSubjects...)

	var resp Response
	var msgErr error

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			if len(subjects) == 0 {
				return
			}

			select {
			case <-time.After(5 * time.Second):
				msgErr = errors.New("timeout")
				return
			case m := <-r.receiver.C():
				if i := slices.Index(subjects, m.Subject); i != -1 {
					resp.Messages = append(resp.Messages, m)
					subjects = append(subjects[:i], subjects[i+1:]...)
				}
				if slices.Contains(errorSubjects, m.Subject) {
					msgErr = NewErrorFromMessage(msg)
					return
				}
			}
		}
	}()

	err := r.sender.SendRaw(msg.Encode())
	if err != nil {
		return nil, fmt.Errorf("send: %w", err)
	}

	wg.Wait()

	if msgErr != nil {
		return nil, fmt.Errorf("send response: %w", msgErr)
	}

	return &resp, nil
}
