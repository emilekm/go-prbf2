package prism

import (
	"bytes"
)

type BufferMessage interface {
	Subject() Subject
}

type ChatMessageType int

const (
	ChatMessageUnknown ChatMessageType = iota - 1
	ChatMessageTypeOpfor
	ChatMessageTypeBlufor
	ChatMessageTypeSquad
	ChatMessageTypeServerMessage
	ChatMessageTypeServer
	ChatMessageTypeResponse
	ChatMessageTypeAdminAlert
)

type ChatMessage struct {
	Type       ChatMessageType
	Timestamp  int
	Channel    string
	PlayerName string
	Content    string
}

func (m ChatMessage) Subject() Subject {
	return SubjectChat
}

type KillMessage struct {
	IsTeamKill   bool
	Timestamp    int
	AttackerName string
	VictimName   string
	Weapon       string
}

func (m KillMessage) Subject() Subject {
	return SubjectKill
}

type BufferReceiver[T BufferMessage] struct {
	typ    Subject
	r      *Receiver
	broker *Broker[T]
}

func NewBufferReceiver[T BufferMessage](r *Receiver) *BufferReceiver[T] {
	receiver := &BufferReceiver[T]{
		r:      r,
		broker: NewBroker[T](),
	}

	receiver.Start()

	return receiver
}

func (r *BufferReceiver[T]) Listen() *Subscriber[T] {
	return r.broker.Subscribe()
}

func (r *BufferReceiver[T]) Start() {
	go func() {
		sub := r.r.Listen()
		for msg := range sub.Channel {
			switch msg.Subject {
			case r.typ:
				chatMessages := bytes.Split(msg.Content, []byte(SeparatorBuffer))
				for _, chatMessageData := range chatMessages {
					var chatMsg T
					err := UnmarshalInto(Message{Content: chatMessageData}, &chatMsg)
					if err != nil {
						continue
					}

					r.broker.Publish(chatMsg)
				}
			}
		}
	}()
}
