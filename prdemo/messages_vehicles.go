package prdemo

type VehicleUpdate struct {
	Flags    uint8 `bin:"flags"`
	ID       int16
	Team     *int8     `bin:"flag=1"`
	Position *Position `bin:"flag=2"`
	Rotation *int16    `bin:"flag=4"`
	Health   *int16    `bin:"flag=8"`
}

type VehiclesUpdate []VehicleUpdate

type VehicleAdd struct {
	ID        int16
	Name      string
	MaxHealth uint16
}

type VehiclesAdd []VehicleAdd

type VehicleDestroyed struct {
	ID            int16
	IsKillerKnown bool
	KillerID      uint8
}
