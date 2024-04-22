package prism

import (
	"fmt"
	"sync"
	"time"
)

type SimpleSubscriber chan *RawMessage

type SimpleBroker struct {
	subjectSubscribers map[Subject][]SimpleSubscriber
	subscribers        []SimpleSubscriber

	mutex sync.Mutex
}

func NewSimpleBroker() *SimpleBroker {
	return &SimpleBroker{
		subjectSubscribers: make(map[Subject][]SimpleSubscriber),
		subscribers:        make([]SimpleSubscriber, 0),
	}
}

func (b *SimpleBroker) Subscribe(subject *Subject) SimpleSubscriber {
	if subject == nil {
		subscriber := make(SimpleSubscriber, 1)

		b.addSubscriber(subscriber)

		return subscriber
	}

	subscriber := make(SimpleSubscriber, 1)
	b.addSubscriberWithSubject(subscriber, *subject)

	return subscriber
}

func (b *SimpleBroker) addSubscriberWithSubject(subscriber SimpleSubscriber, subject Subject) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	subscribers, ok := b.subjectSubscribers[subject]
	if !ok {
		subscribers = make([]SimpleSubscriber, 0)
	}

	subscribers = append(subscribers, subscriber)
	b.subjectSubscribers[subject] = subscribers
}

func (b *SimpleBroker) addSubscriber(subscriber SimpleSubscriber) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.subscribers = append(b.subscribers, subscriber)
}

func (b *SimpleBroker) Unsubscribe(subscriber SimpleSubscriber) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	for i, sub := range b.subscribers {
		if sub == subscriber {
			close(sub)
			b.subscribers = append(b.subscribers[:i], b.subscribers[i+1:]...)
			return
		}
	}

	for subject, subscribers := range b.subjectSubscribers {
		for i, sub := range subscribers {
			if sub == subscriber {
				close(sub)
				b.subjectSubscribers[subject] = append(subscribers[:i], subscribers[i+1:]...)
				return
			}
		}
	}
}

func (b *SimpleBroker) Publish(message *RawMessage) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	timer := time.NewTimer(time.Second)

	for _, sub := range append(b.subjectSubscribers[message.Subject()], b.subscribers...) {
		timer.Reset(time.Second)
		select {
		case sub <- message:
		case <-timer.C:
			fmt.Printf("Subscriber slow. Unsubscribing\n")
			b.Unsubscribe(sub)
		}
	}
}

func (b *SimpleBroker) Close() {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	for _, subscriber := range b.subscribers {
		close(subscriber)
	}

	for _, subscribers := range b.subjectSubscribers {
		for _, subscriber := range subscribers {
			close(subscriber)
		}
	}
}
