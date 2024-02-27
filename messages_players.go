package prdemo

type PlayerAdd struct {
	ID   uint8
	IGN  string
	Hash string
	IP   string
}

type PlayersAdd []PlayerAdd

type PlayerRemove struct {
	ID uint8
}

type PlayerVehicle struct {
	ID         int16
	SeatName   *string
	SeatNumber *int8
}

func (p *PlayerVehicle) Decode(m *Message) error {
	err := m.Decode(&p.ID)
	if err != nil {
		return err
	}

	if p.ID >= 0 {
		err = m.Decode(&p.SeatName)
		if err != nil {
			return err
		}

		err = m.Decode(&p.SeatNumber)
		if err != nil {
			return err
		}
	}

	return nil
}

type PlayerUpdateFlag uint16

const (
	PlayerUpdateFlagTeam PlayerUpdateFlag = 1 << iota
	PlayerUpdateFlagSquad
	PlayerUpdateFlagVehicle
	PlayerUpdateFlagHealth
	PlayerUpdateFlagScore
	PlayerUpdateFlagTeamworkScore
	PlayerUpdateFlagKills
	// Jump 1 bit, since it's unused
	PlayerUpdateFlagDeaths = 1 << iota << 1
	PlayerUpdateFlagPing
	// Jump 1 bit, since it's unused
	PlayerUpdateFlagIsAlive = 1 << iota << 2
	PlayerUpdateFlagIsJoining
	PlayerUpdateFlagPosition
	PlayerUpdateFlagRotation
	PlayerUpdateFlagKitName
)

type PlayerUpdate struct {
	Flags         uint16
	ID            uint8
	Team          int8
	Squad         uint8
	Vehicle       PlayerVehicle
	Health        int8
	Score         int16
	TeamworkScore int16
	Kills         int16
	Deaths        int16
	Ping          int16
	IsAlive       bool
	IsJoining     bool
	Position      Position
	Rotation      int16
	KitName       string
}

func (p *PlayerUpdate) Decode(m *Message) error {
	err := m.Decode(&p.Flags)
	if err != nil {
		return err
	}

	err = m.Decode(&p.ID)
	if err != nil {
		return err
	}

	flagToField := map[PlayerUpdateFlag]interface{}{
		PlayerUpdateFlagTeam:          &p.Team,
		PlayerUpdateFlagSquad:         &p.Squad,
		PlayerUpdateFlagVehicle:       &p.Vehicle,
		PlayerUpdateFlagHealth:        &p.Health,
		PlayerUpdateFlagScore:         &p.Score,
		PlayerUpdateFlagTeamworkScore: &p.TeamworkScore,
		PlayerUpdateFlagKills:         &p.Kills,
		PlayerUpdateFlagDeaths:        &p.Deaths,
		PlayerUpdateFlagPing:          &p.Ping,
		PlayerUpdateFlagIsAlive:       &p.IsAlive,
		PlayerUpdateFlagIsJoining:     &p.IsJoining,
		PlayerUpdateFlagPosition:      &p.Position,
		PlayerUpdateFlagRotation:      &p.Rotation,
		PlayerUpdateFlagKitName:       &p.KitName,
	}

	for _, flag := range sortedKeys(flagToField) {
		if p.Flags&uint16(flag) != 0 {
			err = m.Decode(flagToField[flag])
			if err != nil {
				return err
			}
		}
	}

	return nil
}

type PlayersUpdate []PlayerUpdate

type Kill struct {
	AttackerID uint8
	VictimID   uint8
	Weapon     string
}

type Chat struct {
	Channel  uint8
	PlayerID uint8
	Message  string
}

type Revive struct {
	MedicID   uint8
	RevivedID uint8
}

type KitAllocated struct {
	PlayerID uint8
	KitName  string
}

type ProjAdd struct {
	ID           uint16
	PlayerID     uint8
	Type         uint8
	TemplateName string
	Rotation     int16
	Position     Position
}

type ProjUpdate struct {
	ID       uint16
	Rotation int16
	Position Position
}

type ProjRemove struct {
	ID uint16
}
