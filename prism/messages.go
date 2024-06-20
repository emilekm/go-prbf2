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
	return CommandLogin1
}

type Login1Response struct {
	Hash            []byte
	ServerChallenge []byte
}

type Login2Request struct {
	ChallengeDigest string
}

func (l Login2Request) Subject() Subject {
	return CommandLogin2
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
	chats, err := multipartBody[ChatMessage](content)
	if err != nil {
		return err
	}

	*m = chats
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
	kills, err := multipartBody[KillMessage](content)
	if err != nil {
		return err
	}

	*m = kills
	return nil
}

type Map struct {
	Name  string
	Mode  string
	Layer string
}

type ServerDetails struct {
	Name        string
	IP          string
	Port        string
	StartTime   float64
	RoundWarmup int
	RoundLength int
	MaxPlayers  int

	Status         string
	Map            Map
	RoundStartTime float64
	Players        int
	Team1          string
	Team2          string
	Tickets1       int
	Tickets2       int
	// RCONUsers []string // Custom Unmarshaling needed
}

// Player returned with `listplayers` message
type Player struct {
	// Header
	Name          string
	IsAIPlayer    int
	Hash          string
	IP            string
	ProfileID     string
	Index         int
	JoinTimestamp int

	// Details
	Team          int
	Squad         string
	Kit           string
	Vehicle       string
	Score         int
	ScoreTeamwork int
	Kills         int
	Teamkills     int
	Deaths        int
	Valid         int // ???
	Ping          int
	Idle          int
	Alive         int
	Joining       int
	Position      string
	Rotation      string
}

type Players []Player

func (p *Players) UnmarshalMessage(content []byte) error {
	players, err := multipartBody[Player](content)
	if err != nil {
		return err
	}

	*p = players
	return nil
}

// User returned with `getusers` message
type User struct {
	Name  string
	Power int
}

// List of users returned with `getusers` message
type Users []User

func (u *Users) UnmarshalMessage(content []byte) error {
	users, err := multipartBody[User](content)
	if err != nil {
		return err
	}

	*u = users
	return nil
}

type AddUser struct {
	Name     string
	Password string
	Power    int
}

func (_ AddUser) Subject() Subject {
	return CommandAddUser
}

type ChangeUser struct {
	Name        string
	NewName     string
	NewPassword string
	NewPower    int
}

func (_ ChangeUser) Subject() Subject {
	return CommandChangeUser
}

func multipartBody[T any](content []byte) ([]T, error) {
	messages := bytes.Split(content, SeparatorBuffer)

	var result []T
	for _, message := range messages {
		var msg T
		err := UnmarshalMessage(message, &msg)
		if err != nil {
			return nil, err
		}

		result = append(result, msg)
	}

	return result, nil
}
