package prism

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"
)

type Subscriber chan *Message

type broker struct {
	client *Client

	subjectSubscribers map[Subject]map[Subscriber]struct{}
	subscribers        map[Subscriber]struct{}

	mutex  sync.Mutex
	cancel context.CancelFunc
}

func newBroker(client *Client) *broker {
	return &broker{
		client:             client,
		subjectSubscribers: make(map[Subject]map[Subscriber]struct{}),
		subscribers:        make(map[Subscriber]struct{}),
	}
}

// Subscribe creates a new subscriber and adds it to the broker.
func (b *broker) Subscribe(subject Subject) Subscriber {
	subscriber := make(Subscriber, 1)
	b.addSubscriberWithSubject(subscriber, subject)

	if b.cancel == nil {
		defer b.start()
	}

	return subscriber
}

// SubscribeAll creates a new subscriber that listens to all subjects
// and adds it to the broker.
func (b *broker) SubscribeAll() Subscriber {
	subscriber := make(Subscriber, 1)
	b.addSubscriber(subscriber)

	if b.cancel == nil {
		defer b.start()
	}

	return subscriber
}

func (b *broker) start() {
	ctx, cancel := context.WithCancel(context.Background())
	b.cancel = cancel

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				msg, err := b.client.ReadMessage()
				if err != nil {
					slog.Error("Received error when reading message in broker", "err", err)
					continue
				}

				b.publish(msg)
			}
		}
	}()
}

func (b *broker) addSubscriberWithSubject(subscriber Subscriber, subject Subject) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	subscribers, ok := b.subjectSubscribers[subject]
	if !ok {
		subscribers = make(map[Subscriber]struct{})
	}

	subscribers[subscriber] = struct{}{}
	b.subjectSubscribers[subject] = subscribers
}

func (b *broker) addSubscriber(subscriber Subscriber) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.subscribers[subscriber] = struct{}{}
}

// Unsubscribe removes a subscriber from the broker.
func (b *broker) Unsubscribe(subscriber Subscriber) {
	defer func() {
		if len(b.subscribers) == 0 && len(b.subjectSubscribers) == 0 {
			b.cancel()
			b.cancel = nil
		}
	}()

	b.mutex.Lock()
	defer b.mutex.Unlock()

	if _, ok := b.subscribers[subscriber]; ok {
		close(subscriber)
		delete(b.subscribers, subscriber)
		return
	}

	for subject, subscribers := range b.subjectSubscribers {
		if _, ok := subscribers[subscriber]; ok {
			close(subscriber)
			delete(subscribers, subscriber)
			b.subjectSubscribers[subject] = subscribers
		}
	}
}

func (b *broker) publish(message *Message) {
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

	for sub := range b.subjectSubscribers[message.Subject()] {
		publishFn(sub)
	}

	for sub := range b.subscribers {
		publishFn(sub)
	}
}

// Close closes all subscribers.
func (b *broker) Close() {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	for sub := range b.subscribers {
		close(sub)
	}

	for _, subscribers := range b.subjectSubscribers {
		for sub := range subscribers {
			close(sub)
		}
	}

	if b.cancel != nil {
		b.cancel()
		b.cancel = nil
	}
}
