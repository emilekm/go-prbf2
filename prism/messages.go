package prism

import (
	"bytes"
	"fmt"
)

//go:generate go run golang.org/x/tools/cmd/stringer -type=Layer -linecomment -output=messages_strings.go

type Position struct {
	X, Y, Z float64
}

// Unmarshal from string (-120, 40, -138)
func (p *Position) UnmarshalMessage(content []byte) error {
	_, err := fmt.Sscanf(string(content), "(%f, %f, %f)", &p.X, &p.Y, &p.Z)
	return err
}

type RACommandOutcome struct {
	Topic   string
	Content string
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
	Timestamp  float64
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
	Timestamp    float64
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
	ConnectedUsers ConnectedUsers
}

type ControlPoint struct {
	ID   string
	Team int
}

type Fob struct {
	Position Position
	Team     int
}

type Rally struct {
	Position Position
	Team     int
	Squad    int
}

// Not implemented, empty field
// type Objective struct{}

type GameplayDetails struct {
	ControlPoints []ControlPoint
	Fobs          []Fob
	Rallies       []Rally
	// Objectives    []Objective
}

func (g *GameplayDetails) UnmarshalMessage(content []byte) error {
	fields := bytes.Split(content, SeparatorField)

	if len(fields) < 4 {
		return fmt.Errorf("GameplayDetails: invalid number of fields")
	}

	// ControlPoints
	controlPoints := bytes.Split(fields[0], SeparatorBuffer)

	for _, cp := range controlPoints {
		split := bytes.SplitN(cp, []byte(":"), 2)
		g.ControlPoints = append(g.ControlPoints, ControlPoint{
			ID:   string(split[0]),
			Team: int(split[1][0]),
		})
	}

	// Fobs
	fobs := bytes.Split(fields[1], SeparatorBuffer)

	for _, fob := range fobs {
		split := bytes.SplitN(fob, []byte(":"), 2)
		var pos Position
		err := pos.UnmarshalMessage(split[0])
		if err != nil {
			return err
		}

		g.Fobs = append(g.Fobs, Fob{
			Position: pos,
			Team:     int(split[1][0]),
		})
	}

	// Rallies
	rallies := bytes.Split(fields[2], SeparatorBuffer)

	for _, rally := range rallies {
		split := bytes.SplitN(rally, []byte(":"), 3)
		var pos Position
		err := pos.UnmarshalMessage(split[0])
		if err != nil {
			return err
		}

		g.Rallies = append(g.Rallies, Rally{
			Position: pos,
			Team:     int(split[1][0]),
			Squad:    int(split[2][0]),
		})
	}

	return nil
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
	Position      Position
	Rotation      string
}

type PlayerHeader struct {
	Name          string
	IsAIPlayer    bool
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

const (
	fullPlayerFieldsCount   = 23
	updatePlayerFieldsCount = 18
)

type UpdatePlayers []UpdatePlayer

func (p *UpdatePlayers) UnmarshalMessage(content []byte) error {
	messages := bytes.Split(content, SeparatorBuffer)

	var players []UpdatePlayer

	for _, message := range messages {
		fieldsNum := bytes.Count(message, SeparatorField) + 1
		if fieldsNum == fullPlayerFieldsCount {
			var player FullPlayer
			err := Unmarshal(message, &player)
			if err != nil {
				return err
			}

			players = append(players, UpdatePlayer{Full: &player})
		} else if fieldsNum == updatePlayerFieldsCount {
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
