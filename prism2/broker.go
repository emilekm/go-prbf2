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
func (b *Broker) Subscribe(subject Subject) Subscriber {
	subscriber := make(Subscriber, 1)
	b.addSubscriberWithSubject(subscriber, subject)
	return subscriber
}

// SubscribeAll creates a new subscriber that listens to all subjects
// and adds it to the broker.
func (b *Broker) SubscribeAll() Subscriber {
	subscriber := make(Subscriber, 1)
	b.addSubscriber(subscriber)
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

func (b *Broker) publish(message *Message) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	timer := time.NewTimer(time.Second)

	publishFn := func(sub Subscriber) {
		timer.Reset(time.Second)
		select {
		case sub <- message:
		case <-timer.C:
			fmt.Printf("Subscriber slow. Unsubscribing\n")
			b.Unsubscribe(sub)
		}
	}

	for _, sub := range b.subjectSubscribers[message.Subject()] {
		publishFn(sub)
	}

	for _, sub := range b.subscribers {
		publishFn(sub)
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
