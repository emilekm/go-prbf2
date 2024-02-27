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

func NewPlayerVehicle() *PlayerVehicle {
	return &PlayerVehicle{
		SeatName:   new(string),
		SeatNumber: new(int8),
	}
}

func (p PlayerVehicle) Read(m *Message) (any, error) {
	err := m.Decode(&p.ID)
	if err != nil {
		return p, err
	}

	if p.ID >= 0 {
		p.SeatName = new(string)
		err = m.Decode(p.SeatName)
		if err != nil {
			return p, err
		}

		p.SeatNumber = new(int8)
		err = m.Decode(p.SeatNumber)
		if err != nil {
			return p, err
		}
	}

	return p, nil
}

func (p *PlayerVehicle) Decode(m *Message) error {
	nP, err := p.Read(m)
	if err != nil {
		return err
	}

	*(p) = nP.(PlayerVehicle)
	return nil
}

type PlayerUpdate struct {
	Flags         uint16 `bin:"flags"`
	ID            uint8
	Team          *int8          `bin:"flag=1"`
	Squad         *uint8         `bin:"flag=2"`
	Vehicle       *PlayerVehicle `bin:"flag=4"`
	Health        *int8          `bin:"flag=8"`
	Score         *int16         `bin:"flag=16"`
	TeamworkScore *int16         `bin:"flag=32"`
	Kills         *int16         `bin:"flag=64"`
	Deaths        *int16         `bin:"flag=256"`
	Ping          *int16         `bin:"flag=512"`
	IsAlive       *bool          `bin:"flag=2048"`
	IsJoining     *bool          `bin:"flag=4096"`
	Position      *Position      `bin:"flag=8192"`
	Rotation      *int16         `bin:"flag=16384"`
	KitName       *string        `bin:"flag=32768"`
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
