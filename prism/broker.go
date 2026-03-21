package prism

import (
	"context"
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
	return subscriber
}

// SubscribeAll creates a new subscriber that listens to all subjects
// and adds it to the broker.
func (b *broker) SubscribeAll() Subscriber {
	subscriber := make(Subscriber, 1)
	b.addSubscriber(subscriber)
	return subscriber
}

// startLocked starts the broker read loop. Must be called with b.mutex held.
func (b *broker) startLocked() {
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
					slog.Error("Connection lost", "err", err)
					b.Close()
					return
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

	if b.cancel == nil {
		b.startLocked()
	}
}

func (b *broker) addSubscriber(subscriber Subscriber) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.subscribers[subscriber] = struct{}{}

	if b.cancel == nil {
		b.startLocked()
	}
}

// unsubscribeLocked removes and closes a subscriber. Must be called with b.mutex held.
func (b *broker) unsubscribeLocked(subscriber Subscriber) {
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
			return
		}
	}
}

// Unsubscribe removes a subscriber from the broker.
func (b *broker) Unsubscribe(subscriber Subscriber) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.unsubscribeLocked(subscriber)

	if len(b.subscribers) == 0 && len(b.subjectSubscribers) == 0 && b.cancel != nil {
		b.cancel()
		b.cancel = nil
	}
}

func (b *broker) publish(message *Message) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	publishFn := func(sub Subscriber) {
		timer := time.NewTimer(time.Second)
		select {
		case sub <- message:
			timer.Stop()
		case <-timer.C:
			slog.Warn("subscriber slow, unsubscribing")
			b.unsubscribeLocked(sub)
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
