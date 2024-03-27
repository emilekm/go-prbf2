package prism

import (
	"fmt"
	"sync"
	"time"
)

type Subscriber[T any] struct {
	Channel     chan T
	Unsubscribe chan bool
}

type Broker[T any] struct {
	subscribers []*Subscriber[T]
	mutex       sync.Mutex
}

func NewBroker[T any]() *Broker[T] {
	return &Broker[T]{
		subscribers: make([]*Subscriber[T], 0),
	}
}

func (b *Broker[T]) Subscribe() *Subscriber[T] {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	subscriber := &Subscriber[T]{
		Channel:     make(chan T, 1),
		Unsubscribe: make(chan bool),
	}

	b.subscribers = append(b.subscribers, subscriber)

	return subscriber
}

func (b *Broker[T]) Unsubscribe(subscriber *Subscriber[T]) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	for i, sub := range b.subscribers {
		if sub == subscriber {
			close(sub.Channel)
			b.subscribers = append(b.subscribers[:i], b.subscribers[i+1:]...)
			return
		}
	}
}

func (b *Broker[T]) Publish(payload T) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	for _, sub := range b.subscribers {
		select {
		case sub.Channel <- payload:
		case <-time.After(time.Second):
			fmt.Printf("Subscriber slow. Unsubscribing\n")
			b.Unsubscribe(sub)
		}
	}
}
