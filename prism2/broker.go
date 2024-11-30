package prism2

import (
	"fmt"
	"sync"
	"time"
)

type Subscriber chan *Message

type Broker struct {
	subjectSubscribers map[Subject][]Subscriber
	subscribers        []Subscriber

	mutex sync.Mutex
}

func NewBroker() *Broker {
	return &Broker{
		subjectSubscribers: make(map[Subject][]Subscriber),
		subscribers:        make([]Subscriber, 0),
	}
}

// Subscribe creates a new subscriber and adds it to the broker.
// Pass nil subject to subscribe to all subjects.
func (b *Broker) Subscribe(subject *Subject) Subscriber {
	subscriber := make(Subscriber, 1)

	if subject != nil {
		b.addSubscriberWithSubject(subscriber, *subject)
	} else {
		b.addSubscriber(subscriber)
	}

	return subscriber
}

func (b *Broker) addSubscriberWithSubject(subscriber Subscriber, subject Subject) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	subscribers, ok := b.subjectSubscribers[subject]
	if !ok {
		subscribers = make([]Subscriber, 0)
	}

	subscribers = append(subscribers, subscriber)
	b.subjectSubscribers[subject] = subscribers
}

func (b *Broker) addSubscriber(subscriber Subscriber) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.subscribers = append(b.subscribers, subscriber)
}

// Unsubscribe removes a subscriber from the broker.
func (b *Broker) Unsubscribe(subscriber Subscriber) {
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

// Publish sends a message to all subscribers.
// If subscriber is slow to receive, it will be unsubscribed.
func (b *Broker) Publish(message *Message) {
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

// Close closes all subscribers.
func (b *Broker) Close() {
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
