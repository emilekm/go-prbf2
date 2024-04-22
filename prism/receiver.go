package prism

import (
	"context"
	"log/slog"
)

type Receiver struct {
	r      *Reader
	broker *SimpleBroker
	cancel context.CancelFunc
}

func NewReceiver(r *Reader) *Receiver {
	receiver := &Receiver{
		r:      r,
		broker: NewSimpleBroker(),
	}

	receiver.receive()

	return receiver
}

func (r *Receiver) receive() {
	ctx, cancel := context.WithCancel(context.Background())
	r.cancel = cancel

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				rawMsg, err := r.r.ReadMessage()
				if err != nil {
					slog.Warn("failed to read message", "error", err)
					continue
				}
				r.broker.Publish(rawMsg)
			}
		}
	}()
}

func (r *Receiver) Subscribe(subject *Subject) SimpleSubscriber {
	return r.broker.Subscribe(subject)
}

func (r *Receiver) Unsubscribe(sub SimpleSubscriber) {
	r.broker.Unsubscribe(sub)
}

func (r *Receiver) Close() {
	r.broker.Close()
	r.cancel()
}
