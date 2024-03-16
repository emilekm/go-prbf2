package prdemo

type FobAdd struct {
	ID       int32
	Team     uint8
	Position Position
}

type FobsAdd struct {
	Fobs []FobAdd
}

type FobRemove struct {
	ID int32
}

type FobsRemove []FobRemove

type RallyAdd struct {
	TeamSquad uint8
	Position  Position
}

type RallyRemove struct {
	TeamSquad uint8
}

type Tickets struct {
	Tickets int16
}

type SquadName struct {
	TeamSquad uint8
	SquadName string
}

type SLOrder struct {
	TeamSquad uint8
	OrderType uint8
	Position  Position
}

type SLOrders []SLOrder
