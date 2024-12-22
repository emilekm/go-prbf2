package prism

import (
	"bytes"
	"reflect"
)

//go:generate go run golang.org/x/tools/cmd/stringer -type=Layer -linecomment -output=messages_strings.go

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
	chats, err := UnmarshalMultipartBody[ChatMessage](content)
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
	kills, err := UnmarshalMultipartBody[KillMessage](content)
	if err != nil {
		return err
	}

	*m = kills
	return nil
}

type Layer int

const (
	LayerUnknown     Layer = 0   // ???
	LayerInfantry    Layer = 16  // inf
	LayerAlternative Layer = 32  // alt
	LayerStandard    Layer = 64  // std
	LayerLarge       Layer = 128 // lrg
)

type Map struct {
	Name  string
	Mode  string
	Layer Layer
}

type ConnectedUsers []string

func (c *ConnectedUsers) UnmarshalMessage(content []byte) error {
	users := bytes.Split(content, SeparatorBuffer)
	for _, user := range users {
		*c = append(*c, string(user))
	}

	return nil
}

// Subjects:
// - serverdetails
// - updateserverdetails
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
	// []string separated by SeparatorBuffer
	// Currently impossible to unmarshal this field
	ConnectedUsers ConnectedUsers
}

type PlayerDetails struct {
	Team int
	// 0 = none, ?L = squad leader, C = commander
	Squad         string
	Kit           string
	Vehicle       string
	Score         int
	ScoreTeamwork int
	Kills         int
	Teamkills     int
	Deaths        int
	Valid         bool // ???
	Ping          int
	Idle          bool
	Alive         bool
	Joining       bool
	Position      string
	Rotation      string
}

type PlayerHeader struct {
	Name          string
	IsAIPlayer    int
	Hash          string
	IP            string
	ProfileID     string
	Index         int
	JoinTimestamp float64
}

// FullPlayer returned with `listplayers` message
type FullPlayer struct {
	PlayerHeader
	PlayerDetails
}

// Subjects:
// - listplayers
type Players []FullPlayer

func (p *Players) UnmarshalMessage(content []byte) error {
	players, err := UnmarshalMultipartBody[FullPlayer](content)
	if err != nil {
		return err
	}

	*p = players
	return nil
}

type Player struct {
	Name  string
	Index int
	PlayerDetails
}

// Subjects:
// - updateplayers
// NOTE: some players in `updateplayers` might have body of FullPlayer instead of Player
type UpdatePlayer struct {
	Full   *FullPlayer
	Update *Player
}

type UpdatePlayers []UpdatePlayer

func (p *UpdatePlayers) UnmarshalMessage(content []byte) error {
	messages := bytes.Split(content, SeparatorBuffer)

	var players []UpdatePlayer

	for _, message := range messages {
		fieldsNum := len(bytes.Split(message, SeparatorField))
		if fieldsNum == reflect.TypeOf(FullPlayer{}).NumField() {
			var player FullPlayer
			err := Unmarshal(message, &player)
			if err != nil {
				return err
			}

			players = append(players, UpdatePlayer{Full: &player})
		} else if fieldsNum == reflect.TypeOf(Player{}).NumField() {
			var player Player
			err := Unmarshal(message, &player)
			if err != nil {
				return err
			}

			players = append(players, UpdatePlayer{Update: &player})
		}
	}

	*p = players
	return nil
}

func UnmarshalMultipartBody[T any](content []byte) ([]T, error) {
	messages := bytes.Split(content, SeparatorBuffer)

	var result []T
	for _, message := range messages {
		var msg T
		err := Unmarshal(message, &msg)
		if err != nil {
			return nil, err
		}

		result = append(result, msg)
	}

	return result, nil
}
