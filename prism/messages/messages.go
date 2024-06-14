package messages

import (
	"bytes"

	"github.com/emilekm/go-prbf2/prism"
)

type baseMessage struct{}

type Login1Request struct {
	baseMessage
	ServerVersion      prism.ServerVersion
	Username           string
	ClientChallengeKey []byte
}

func (l Login1Request) Subject() prism.Subject {
	return prism.SubjectLogin1
}

type Login1Response struct {
	baseMessage
	Hash            []byte
	ServerChallenge []byte
}

type Login2Request struct {
	baseMessage
	ChallengeDigest string
}

func (l Login2Request) Subject() prism.Subject {
	return prism.SubjectLogin2
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
	baseMessage
	Type       ChatMessageType
	Timestamp  int
	Channel    string
	PlayerName string
	Content    string
}

type ChatMessages []ChatMessage

func (m *ChatMessages) Decode(content []byte) error {
	messages := bytes.Split(content, prism.SeparatorBuffer)

	for _, message := range messages {
		var msg ChatMessage
		err := msg.Decode(message)
		if err != nil {
			return err
		}

		*m = append(*m, msg)
	}

	return nil
}

type KillMessage struct {
	baseMessage
	IsTeamKill   bool
	Timestamp    int
	AttackerName string
	VictimName   string
	Weapon       string
}

type KillMessages []KillMessage

func (m *KillMessages) Decode(content []byte) error {
	messages := bytes.Split(content, prism.SeparatorBuffer)

	for _, message := range messages {
		var msg KillMessage
		err := msg.Decode(message)
		if err != nil {
			return err
		}

		*m = append(*m, msg)
	}

	return nil
}
