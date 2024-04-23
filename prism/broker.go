package prism

import (
	"context"
	"fmt"
	"sync"
	"time"
)

var emptySubject = Subject("")

type Subscriber[T any] chan T

type Broker[T Message] struct {
	subscribers map[Subject][]Subscriber[T]
	mutex       sync.Mutex
}

func NewBroker[T Message]() *Broker[T] {
	return &Broker[T]{
		subscribers: make(map[Subject][]Subscriber[T]),
	}
}

func (b *Broker[T]) Subscribe(subject *Subject) Subscriber[T] {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	subscriber := make(Subscriber[T], 1)

	if subject == nil {
		subject = &emptySubject
	}

	subscribers, ok := b.subscribers[*subject]
	if !ok {
		subscribers = make([]Subscriber[T], 0)
	}

	subscribers = append(subscribers, subscriber)

	return subscriber
}

func (b *Broker[T]) Unsubscribe(subscriber Subscriber[T]) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	for _, subscribers := range b.subscribers {
		for i, sub := range subscribers {
			if sub == subscriber {
				close(sub)
				subscribers = append(subscribers[:i], subscribers[i+1:]...)
				return
			}
		}
	}
}

func (b *Broker[T]) Publish(payload T) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.publishSubject(payload.Subject(), payload)
	b.publishSubject(emptySubject, payload)
}

func (b *Broker[T]) publishSubject(subject Subject, payload T) {
	timer := time.NewTimer(time.Second)

	for _, sub := range b.subscribers[subject] {
		timer.Reset(time.Second)
		select {
		case sub <- payload:
		case <-timer.C:
			fmt.Printf("Subscriber slow. Unsubscribing\n")
			b.Unsubscribe(sub)
		}
	}

	timer.Stop()
}

type PassThroughBroker[T Message] struct {
	Broker[T]
	mainBroker *Broker[RawMessage]

	subject Subject
	cancel  context.CancelFunc
}

func NewPassThroughBroker[T Message](mainBroker *Broker[RawMessage]) *PassThroughBroker[T] {
	return &PassThroughBroker[T]{
		mainBroker: mainBroker,
		subject:    (*new(T)).Subject(),
	}
}

func (b *PassThroughBroker[T]) Subscribe() Subscriber[T] {
	if len(b.subscribers[b.subject]) == 0 {
		b.enablePassThrough()
	}
	return b.Broker.Subscribe(&b.subject)
}

func (b *PassThroughBroker[T]) Unsubscribe(subscriber Subscriber[T]) {
	b.Broker.Unsubscribe(subscriber)
	if len(b.subscribers[b.subject]) == 0 {
		b.cancel()
	}
}

func (b *PassThroughBroker[T]) enablePassThrough() {
	ctx, cancel := context.WithCancel(context.Background())
	b.cancel = cancel

	sub := b.mainBroker.Subscribe(&b.subject)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case rawMsg := <-sub:
				if rawMsg.Subject() != b.subject {
					continue
				}

				var msg T
				err := decodeContent(rawMsg.Content(), &msg)
				if err != nil {
					fmt.Printf("Failed to decode message: %v\n", err)
					continue
				}
				b.Publish(msg)
			}
		}
	}()
}

// func findBroker[T Message](brokerMap map[string]any, mainBroker *Broker[RawMessage]) *PassThroughBroker[T] {
// 	var msg T
// 	topic := fmt.Sprintf("%s.%T", msg.Subject(), msg)
//
// 	if _, ok := brokerMap[topic]; !ok {
// 		brokerMap[topic] = NewPassThroughBroker[T](mainBroker)
// 	}
//
// 	return brokerMap[topic].(*PassThroughBroker[T])
// }
//
// func SubscribeToSpecificBroker[T Message](r *Receiver) Subscriber[T] {
// 	b := findBroker[T](r.specificBrokers, r.mainBroker)
//
// 	return b.Subscribe()
// }
