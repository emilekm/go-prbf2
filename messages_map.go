package prdemo

import "time"

type Map struct {
	Name     string
	Gamemode string
	Layer    uint8
}

type ServerDetails struct {
	Version         int32
	DemoTimePerTick float32
	IPPort          string
	ServerName      string
	MaxPlayers      uint8
	RoundLength     uint16
	BriefingTime    uint16
	Map             Map
	BluforTeam      string
	OpforTeam       string
	StartTime       time.Time
	Tickets1        uint16
	Tickets2        uint16
	MapSize         float32
}

type Position struct {
	X int16
	Y int16
	Z int16
}

type CacheAdd struct {
	ID       uint8
	Position Position
}

type CachesAdd []CacheAdd

type CacheRemove struct {
	ID uint8
}

type CacheReveal struct {
	ID uint8
}

type CachesReveal []CacheReveal

type IntelChange struct {
	IntelCount int8
}

type DoD struct {
	Team        uint8
	Inverted    uint8
	DoDType     uint8
	NumPoints   uint8
	PointsArray []struct {
		X float32
		Y float32
	}
}

type DoDs []DoD

type Flag struct {
	CpID       int16
	OwningTeam uint8
	Position   Position
	Radius     uint16
}

type Flags []Flag
