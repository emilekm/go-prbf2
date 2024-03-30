package prism

import (
	"context"
	"errors"
	"slices"
	"sync"
	"time"
)

const defaultTimeout = 5 * time.Second

type Responder struct {
	c       *Client
	Timeout time.Duration
}

func NewResponder(c *Client) *Responder {
	return &Responder{
		c:       c,
		Timeout: defaultTimeout,
	}
}

type ResponseOption func(*responseOptions)

type responseFilter func(Message) (Message, bool)

type responseOptions struct {
	filters []responseFilter
}

func filterOnChatMessage(m Message, typ ChatMessageType) (Message, bool) {
	if m.Subject() != SubjectChat {
		return nil, false
	}

	chatMsgs, ok := m.(ChatMessages)
	if !ok {
		return nil, false
	}

	for _, chatMsg := range chatMsgs {
		if chatMsg.Type == typ {
			return chatMsg, true
		}
	}

	return nil, false
}

func ResponseWithChatMessage(typ ChatMessageType) ResponseOption {
	return func(o *responseOptions) {
		o.filters = append(o.filters, func(m Message) (Message, bool) {
			return filterOnChatMessage(m, typ)
		})
	}
}

func ResponseWithMessageSubject(sub Subject) ResponseOption {
	return func(o *responseOptions) {
		o.filters = append(o.filters, func(m Message) (Message, bool) {
			if m.Subject() != sub {
				return nil, false
			}
			return m, true
		})
	}
}

type Response struct {
	Messages []Message
}

func (r *Responder) SendWithResponse(ctx context.Context, msg Message, responseOpts ...ResponseOption) (*Response, error) {
	var opts responseOptions
	for _, opt := range responseOpts {
		opt(&opts)
	}

	var resp Response
	var respErr error

	id := r.c.Next()
	r.c.StartRequest(id)

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()

		sub := r.c.Subscribe()
		defer r.c.Unsubscribe(sub)

		timer := time.NewTimer(r.c.Timeout)

		filters := opts.filters[:]
		for {
			if len(filters) == 0 {
				return
			}

			select {
			case <-ctx.Done():
				respErr = ctx.Err()
				return
			case <-timer.C:
				respErr = errors.New("timeout")
				return
			case m := <-sub:
				if isErrorMessage(m) {
					respErr = m.(error)
					return
				}

				if msg, ok := filters[0](m); ok {
					resp.Messages = append(resp.Messages, msg)
					if len(filters) > 0 {
						filters = filters[1:]
					}
				}
			}
		}
	}()

	if err := r.c.Send(msg); err != nil {
		return nil, err
	}

	wg.Wait()

	r.c.EndRequest(id)

	if respErr != nil {
		return nil, respErr
	}

	return &resp, nil
}

func isErrorMessage(m Message) bool {
	return slices.Contains(errorSubjects, m.Subject())
}
