package prdemo

type VehicleUpdateFlag uint8

const (
	VehicleUpdateFlagTeam VehicleUpdateFlag = 1 << iota
	VehicleUpdateFlagPosition
	VehicleUpdateFlagRotation
	VehicleUpdateFlagHealth
)

type VehicleUpdate struct {
	Flags    uint8
	ID       int16
	Team     *int8
	Position *Position
	Rotation *int16
	Health   *int8
}

func (v *VehicleUpdate) Decode(d DemoReader) error {
	err := d.Decode(&v.Flags)
	if err != nil {
		return err
	}

	err = d.Decode(&v.ID)
	if err != nil {
		return err
	}

	flagToField := map[VehicleUpdateFlag]interface{}{
		VehicleUpdateFlagTeam:     v.Team,
		VehicleUpdateFlagPosition: v.Position,
		VehicleUpdateFlagRotation: v.Rotation,
		VehicleUpdateFlagHealth:   v.Health,
	}

	for _, flag := range sortedKeys(flagToField) {
		if v.Flags&uint8(flag) != 0 {
			err = d.Decode(flagToField[flag])
			if err != nil {
				return err
			}
		}
	}

	return nil
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
