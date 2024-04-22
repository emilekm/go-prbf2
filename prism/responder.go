package prism

import (
	"context"
	"net/textproto"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

const defaultTimeout = 5 * time.Second

type ResponseOption func(*responseOptions)

type responseFilter struct {
	Subject *Subject
	F       func(Message) (Message, bool)
}

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
		s := SubjectChat
		o.filters = append(o.filters, responseFilter{
			Subject: &s,
			F: func(m Message) (Message, bool) {
				return filterOnChatMessage(m, typ)
			},
		})
	}
}

func ResponseWithMessageSubject(sub Subject) ResponseOption {
	return func(o *responseOptions) {
		o.filters = append(o.filters, responseFilter{
			Subject: &sub,
		})
	}
}

type Response struct {
	Messages []Message
}

type Responder struct {
	textproto.Pipeline
	receiver *Receiver
	writer   *Writer
	Timeout  time.Duration
}

func NewResponder(receiver *Receiver, writer *Writer) *Responder {
	return &Responder{
		receiver: receiver,
		writer:   writer,
		Timeout:  defaultTimeout,
	}
}

func (r *Responder) SendWithResponse(ctx context.Context, msg Message, responseOpts ...ResponseOption) (*Response, error) {
	var opts responseOptions
	for _, opt := range responseOpts {
		opt(&opts)
	}

	var resp Response

	id := r.Next()
	r.StartRequest(id)

	ctx, cancel := context.WithTimeout(ctx, r.Timeout)

	wg := &sync.WaitGroup{}

	errGroup, eCtx := errgroup.WithContext(ctx)

	wg.Add(2)
	errGroup.Go(errorListener(eCtx, r.receiver, SubjectError, wg))
	errGroup.Go(errorListener(eCtx, r.receiver, SubjectCriticalError, wg))

	respCh := make(chan Message, len(opts.filters))

	successGroup, sCtx := errgroup.WithContext(ctx)

	for _, filter := range opts.filters {
		wg.Add(1)
		successGroup.Go(filterListener(sCtx, r.receiver, filter, respCh, wg))
	}

	// Wait until everyone is listening
	wg.Wait()

	if err := r.writer.WriteMessage(msg); err != nil {
		cancel()
		return nil, err
	}

	multiErrGroup, _ := errgroup.WithContext(ctx)
	for _, g := range []*errgroup.Group{errGroup, successGroup} {
		g := g
		multiErrGroup.Go(func() error {
			err := g.Wait()
			cancel()
			return err
		})
	}

	err := multiErrGroup.Wait()
	cancel()
	if err != nil {
		return nil, err
	}

	r.EndRequest(id)

	for msg := range respCh {
		resp.Messages = append(resp.Messages, msg)
		if len(resp.Messages) == len(opts.filters) {
			break
		}
	}

	return &resp, nil
}

func errorListener(ctx context.Context, r *Receiver, errSubject Subject, wg *sync.WaitGroup) func() error {
	return func() error {
		sub := r.Subscribe(&errSubject)
		defer r.Unsubscribe(sub)

		// Mark that we are listening
		wg.Done()

		for {
			select {
			case <-ctx.Done():
				return nil
			case msg := <-sub:
				var respErr Error
				err := DecodeContent(msg.Content(), &respErr)
				if err != nil {
					return err
				}

				return respErr
			}
		}
	}
}

func filterListener(ctx context.Context, r *Receiver, filter responseFilter, c chan Message, wg *sync.WaitGroup) func() error {
	return func() error {
		var sub SimpleSubscriber
		if filter.Subject != nil {
			sub = r.Subscribe(filter.Subject)
		} else {
			sub = r.Subscribe(nil)
		}
		defer r.Unsubscribe(sub)

		// Mark that we are listening
		wg.Done()

		for {
			select {
			case <-ctx.Done():
				return nil
			case msg := <-sub:
				var selectedMsg Message

				if filter.F != nil {
					var ok bool
					selectedMsg, ok = filter.F(msg)
					if !ok {
						continue
					}
				}

				if selectedMsg == nil {
					selectedMsg = msg
				}

				c <- msg
				return nil
			}
		}
	}
}
