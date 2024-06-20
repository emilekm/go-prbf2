package prism

import (
	"bytes"
)

type Login1Request struct {
	ServerVersion      ServerVersion
	Username           string
	ClientChallengeKey []byte
}

func (l Login1Request) Subject() Subject {
	return SubjectLogin1
}

type Login1Response struct {
	Hash            []byte
	ServerChallenge []byte
}

type Login2Request struct {
	ChallengeDigest string
}

func (l Login2Request) Subject() Subject {
	return SubjectLogin2
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

type ChatMessages []ChatMessage

func (m *ChatMessages) UnmarshalMessage(content []byte) error {
	messages := bytes.Split(content, SeparatorBuffer)

	for _, message := range messages {
		var msg ChatMessage
		err := UnmarshalMessage(message, &msg)
		if err != nil {
			return err
		}

		*m = append(*m, msg)
	}

	return nil
}

type KillMessage struct {
	IsTeamKill   bool
	Timestamp    int
	AttackerName string
	VictimName   string
	Weapon       string
}

type KillMessages []KillMessage

func (m *KillMessages) UnmarshalMessage(content []byte) error {
	messages := bytes.Split(content, SeparatorBuffer)

	for _, message := range messages {
		var msg KillMessage
		err := UnmarshalMessage(message, &msg)
		if err != nil {
			return err
		}

		*m = append(*m, msg)
	}

	return nil
}
