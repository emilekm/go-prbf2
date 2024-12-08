package prism

const (
	SeparatorField1 = byte(0x03)
	SeparatorNull1  = byte(0x00)
)

var (
	SeparatorStart   = []byte{0x01}
	SeparatorSubject = []byte{0x02}
	SeparatorField   = []byte{SeparatorField1}
	SeparatorEnd     = []byte{0x04}
	SeparatorNull    = []byte{SeparatorNull1}
	SeparatorBuffer  = []byte{0x0A}
)

// Incoming message subject
type Subject string

const (
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
	CommandListPlayers         Command = "listplayers"
	CommandServerDetails       Command = "serverdetails"
	CommandGameplayDetails     Command = "gameplaydetails"
	CommandReadMaplist         Command = "readmaplist"
	CommandAPIAdmin            Command = "apiadmin"
	CommandServerDetailsAlways Command = "serverdetailsalways"
)

type ServerVersion string

const ServerVersion1 ServerVersion = "1"
