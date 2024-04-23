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
	Subject Subject
}

type responseOptions struct {
	filters []responseFilter
}

func ResponseWithMessageSubject(sub Subject) ResponseOption {
	return func(o *responseOptions) {
		o.filters = append(o.filters, responseFilter{
			Subject: sub,
		})
	}
}

type Response struct {
	Messages []*RawMessage
}

type Responder struct {
	pipeline *textproto.Pipeline
	receiver *Receiver
	writer   *Writer
	Timeout  time.Duration
}

func NewResponder(receiver *Receiver, writer *Writer, pipeline *textproto.Pipeline) *Responder {
	return &Responder{
		receiver: receiver,
		writer:   writer,
		pipeline: pipeline,
		Timeout:  defaultTimeout,
	}
}

func (r *Responder) Send(ctx context.Context, msg Message, responseOpts ...ResponseOption) (*Response, error) {
	var opts responseOptions
	for _, opt := range responseOpts {
		opt(&opts)
	}

	var resp Response

	id := r.pipeline.Next()
	r.pipeline.StartRequest(id)

	ctx, cancel := context.WithTimeout(ctx, r.Timeout)

	wg := &sync.WaitGroup{}

	errGroup, eCtx := errgroup.WithContext(ctx)

	wg.Add(2)
	errGroup.Go(errorListener(eCtx, r.receiver, SubjectError, wg))
	errGroup.Go(errorListener(eCtx, r.receiver, SubjectCriticalError, wg))

	respCh := make(chan *RawMessage, len(opts.filters))

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

	r.pipeline.EndRequest(id)

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
				err := decodeContent(msg.Content(), &respErr)
				if err != nil {
					return err
				}

				return respErr
			}
		}
	}
}

func filterListener(ctx context.Context, r *Receiver, filter responseFilter, c chan *RawMessage, wg *sync.WaitGroup) func() error {
	return func() error {
		sub := r.Subscribe(&filter.Subject)
		defer r.Unsubscribe(sub)

		// Mark that we are listening
		wg.Done()

		for {
			select {
			case <-ctx.Done():
				return nil
			case msg := <-sub:
				c <- msg
				return nil
			}
		}
	}
}
