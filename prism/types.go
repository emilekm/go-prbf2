package prism

const (
	SeparatorField1 = byte(0x03)
	SeparatorNull1  = byte(0x00)
)

var (
	SeparatorStart   = []byte{0x01}
	SeparatorSubject = []byte{0x02}
	SeparatorField   = []byte{0x03}
	SeparatorEnd     = []byte{0x04}
	SeparatorNull    = []byte{0x00}
	SeparatorBuffer  = []byte{0x0A}
)

// Incoming message subject
type Subject string

const (
	SubjectLogin1              Subject = "login1"
	SubjectConnected           Subject = "connected"
	SubjectServerDetails       Subject = "serverdetails"
	SubjectUpdateServerDetails Subject = "updateserverdetails"
	SubjectGameplayDetails     Subject = "gameplaydetails"
	SubjectRAConfig            Subject = "raconfig"
	SubjectMaplist             Subject = "maplist"
	SubjectSuccess             Subject = "success"
	SubjectError               Subject = "error"
	SubjectCriticalError       Subject = "errorcritical"
	SubjectAPIAdminResult      Subject = "APIAdminResult"
	SubjectListPlayers         Subject = "listplayers"
	SubjectUpdatePlayers       Subject = "updateplayers"
	SubjectPlayerLeave         Subject = "playerleave"
	SubjectChat                Subject = "chat"
	SubjectKill                Subject = "kill"
	SubjectRACommand           Subject = "say"
)

// Command is the outgoing subject
type Command = Subject

const (
	CommandLogin1              Command = "login1"
	CommandLogin2              Command = "login2"
	CommandListPlayers         Command = "listplayers"
	CommandServerDetails       Command = "serverdetails"
	CommandGameplayDetails     Command = "gameplaydetails"
	CommandReadMaplist         Command = "readmaplist"
	CommandAPIAdmin            Command = "apiadmin"
	CommandServerDetailsAlways Command = "serverdetailsalways"
)

type ServerVersion string

const ServerVersion1 ServerVersion = "1"
