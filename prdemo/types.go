package prdemo

//go:generate go run golang.org/x/tools/cmd/stringer -type=MessageType -output=types_strings.go

type MessageType uint8

const (
	ServerDetailsType MessageType = 0x00
	DodListType       MessageType = 0x01

	PlayerUpdateType MessageType = 0x10
	PlayerAddType    MessageType = 0x11
	PlayerRemoveType MessageType = 0x12

	VehicleUpdateType    MessageType = 0x20
	VehicleAddType       MessageType = 0x21
	VehicleDestroyedType MessageType = 0x22

	FobAddType    MessageType = 0x30
	FobRemoveType MessageType = 0x31

	FlagUpdateType MessageType = 0x40
	FlagListType   MessageType = 0x41

	KillType MessageType = 0x50
	ChatType MessageType = 0x51

	TicketsTeam1Type MessageType = 0x52
	TicketsTeam2Type MessageType = 0x53

	RallyAddType    MessageType = 0x60
	RallyRemoveType MessageType = 0x61

	CacheAddType    MessageType = 0x70
	CacheRemoveType MessageType = 0x71
	CacheRevealType MessageType = 0x72
	IntelChangeType MessageType = 0x73

	MarkerAddType    MessageType = 0x80
	MarkerRemoveType MessageType = 0x81

	ProjUpdateType MessageType = 0x90
	ProjAddType    MessageType = 0x91
	ProjRemoveType MessageType = 0x92

	ReviveType       MessageType = 0xA0
	KitAllocatedType MessageType = 0xA1
	SquadNameType    MessageType = 0xA2
	SlOrdersType     MessageType = 0xA3

	RoundEndType MessageType = 0xF0
	TicksType    MessageType = 0xF1

	PrivateMessageType MessageType = 0xFD
	ErrorMessageType   MessageType = 0xFE
)
