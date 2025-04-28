package gamespy3

type FieldType byte

const (
	FieldTypeHostName    FieldType = 0x01
	FieldTypeGameName    FieldType = 0x02
	FieldTypeGameVersion FieldType = 0x03
	FieldTypeHostPort    FieldType = 0x04
	FieldTypeMapName     FieldType = 0x05
	FieldTypeGameType    FieldType = 0x06
	FieldTypeGameVariant FieldType = 0x07
	FieldTypeNumPlayers  FieldType = 0x08
	FieldTypeMaxPlayers  FieldType = 0x0a
	FieldTypeGameMode    FieldType = 0x0b
	FieldTypeTimeLimit   FieldType = 0x10
	FieldTypeRoundTime   FieldType = 0x11
	FieldTypePassword    FieldType = 0x13
	FieldTypeDedicated   FieldType = 0x33
	FieldTypeRanked      FieldType = 0x34
	FieldTypeAntiCheat   FieldType = 0x35
	FieldTypeOS          FieldType = 0x36
	FieldTypeUnknown1    FieldType = 0x37
	FieldTypeUnknown2    FieldType = 0x3a
	FieldTypeUnknown3    FieldType = 0x47
	FieldTypeUnknown4    FieldType = 0x48
	FieldTypeUnknown5    FieldType = 0x49
)

type Map struct {
	Name     string `mapstructure:"mapname"`
	GameMode string `mapstructure:"gametype"`
	Layer    int    `mapstructure:"bf2_mapsize"`
}

type Header struct {
	Hostname         string  `mapstructure:"hostname"`
	GameName         string  `mapstructure:"gamename"`
	GameVersion      string  `mapstructure:"gamever"`
	Map              Map     `mapstructure:",squash"`
	GameVariant      string  `mapstructure:"gamevariant"`
	NumPlayers       int     `mapstructure:"numplayers"`
	MaxPlayers       int     `mapstructure:"maxplayers"`
	GameMode         string  `mapstructure:"gamemode"`
	Password         int     `mapstructure:"password"`
	TimeLimit        int     `mapstructure:"timelimit"`
	RoundTime        int     `mapstructure:"roundtime"`
	HostPort         int     `mapstructure:"hostport"`
	Dedicated        int     `mapstructure:"bf2_dedicated"`
	Ranked           int     `mapstructure:"bf2_ranked"`
	AntiCheat        int     `mapstructure:"bf2_anticheat"`
	OS               string  `mapstructure:"bf2_os"`
	AutoRec          int     `mapstructure:"bf2_autorec"`
	DownloadIndex    string  `mapstructure:"bf2_d_idx"`
	DownloadURL      string  `mapstructure:"bf2_d_dl"`
	VoIP             int     `mapstructure:"bf2_voip"`
	AutoBalanced     int     `mapstructure:"bf2_autobalanced"`
	FriendlyFire     int     `mapstructure:"bf2_friendlyfire"`
	TKMode           string  `mapstructure:"bf2_tkmode"`
	StartDelay       int     `mapstructure:"bf2_startdelay"`
	SpawnTime        float64 `mapstructure:"bf2_spawntime"`
	SponsorText      string  `mapstructure:"bf2_sponsortext"`
	SponsorLogoURL   string  `mapstructure:"bf2_sponsorlogo_url"`
	CommunityLogoURL string  `mapstructure:"bf2_communitylogo_url"`
	ScoreLimit       int     `mapstructure:"bf2_scorelimit"`
	TicketRatio      int     `mapstructure:"bf2_ticketratio"`
	TeamRatio        float64 `mapstructure:"bf2_teamratio"`
	Team1            string  `mapstructure:"bf2_team1"`
	Team2            string  `mapstructure:"bf2_team2"`
	Bots             int     `mapstructure:"bf2_bots"`
	Pure             int     `mapstructure:"bf2_pure"`
	GlobalUnlocks    int     `mapstructure:"bf2_globalunlocks"`
	FPS              float64 `mapstructure:"bf2_fps"`
	Plasma           int     `mapstructure:"bf2_plasma"`
	ReservedSlots    int     `mapstructure:"bf2_reservedslots"`
	CoopBotRatio     string  `mapstructure:"bf2_coopbotratio"`
	CoopBotCount     string  `mapstructure:"bf2_coopbotcount"`
	CoopBotDiff      string  `mapstructure:"bf2_coopbotdiff"`
	NoVehicles       int     `mapstructure:"bf2_novehicles"`
}

type Player struct {
	ProfileID string
	Player    string
	Team      string
	Deaths    string
	Score     string
	Ping      string
	Skill     string
	Bot       bool
}

type Team struct {
}

type Teams struct {
	Team1 Team
	Team2 Team
}

type ServerInformation struct {
	Header  Header
	Players []Player `mapstructure:",squash"`
	Teams   Teams    `mapstructure:",squash"`
}
