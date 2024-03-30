package prism

import (
	"fmt"
	"sync"
	"time"
)

type Subscriber[T any] chan T

type Broker[T any] struct {
	subscribers []Subscriber[T]
	mutex       sync.Mutex
}

func NewBroker[T any]() *Broker[T] {
	return &Broker[T]{
		subscribers: make([]Subscriber[T], 0),
	}
}

func (b *Broker[T]) Subscribe() Subscriber[T] {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	subscriber := make(Subscriber[T], 1)

	b.subscribers = append(b.subscribers, subscriber)

	return subscriber
}

func (b *Broker[T]) Unsubscribe(subscriber Subscriber[T]) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	for i, sub := range b.subscribers {
		if sub == subscriber {
			close(sub)
			b.subscribers = append(b.subscribers[:i], b.subscribers[i+1:]...)
			return
		}
	}
}

func (b *Broker[T]) Publish(payload T) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	timer := time.NewTimer(time.Second)

	for _, sub := range b.subscribers {
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
