package prism

import (
	"bytes"
	"fmt"
)

type MessageDecoder interface {
	Decode(RawMessage) error
}

type ServerVersion string

const ServerVersion1 ServerVersion = "1"

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
	Test            string
}

func (l Login1Response) Subject() Subject {
	return SubjectLogin1
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

func (m ChatMessage) Subject() Subject {
	return SubjectChat
}

type ChatMessages []ChatMessage

func (m ChatMessages) Subject() Subject {
	return SubjectChat
}

func (m *ChatMessages) Decode(content []byte) error {
	messages := bytes.Split(content, SeparatorBuffer)

	for _, message := range messages {
		var msg ChatMessage
		err := decodeContent(message, &msg)
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

func (m KillMessage) Subject() Subject {
	return SubjectKill
}

type KillMessages []KillMessage

func (m KillMessages) Subject() Subject {
	return SubjectKill
}

func (m *KillMessages) Decode(content []byte) error {
	messages := bytes.Split(content, SeparatorBuffer)

	for _, message := range messages {
		var msg KillMessage
		err := decodeContent(message, &msg)
		if err != nil {
			return err
		}

		*m = append(*m, msg)
	}

	return nil
}

type ErrorCode int

const (
	ErrorCodeUnknown         ErrorCode = iota
	ErrorCodeUnauthenticated ErrorCode = iota + 3000
	ErrorCodeIncorectUsernameOrPassword
	ErrorCodeInssuficientPermissions
	ErrorCodeAccountExists
	ErrorCodeOwnAccont
	ErrorCodeSuperUserLastAccount
	ErrorCodeDeletedUser
	ErrorCodeServerVersion
)

var errorSubjects = []Subject{
	SubjectCriticalError,
	SubjectError,
}

type Error struct {
	Code    ErrorCode
	Content string
}

func (e Error) Subject() Subject {
	return SubjectError
}

func (e Error) Error() string {
	return fmt.Sprintf("ErrorCode: %d, Content: %s", e.Code, e.Content)
}

type CriticalError Error

func (e CriticalError) Subject() Subject {
	return SubjectCriticalError
}
