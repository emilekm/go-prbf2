package gamespy3

import (
	"bytes"
	"context"
	"net"
	"time"

	"github.com/go-viper/mapstructure/v2"
)

var (
	delim = []byte{0x00, 0x00, 0x01}
	query = []byte{0xFE, 0xFD, 0x00, 0x04, 0x05, 0x06, 0x07, 0xFF, 0x00, 0x00, 0x01}
)

type Map struct {
	Name     string `mapstructure:"mapname"`
	GameMode string `mapstructure:"gametype"`
	Layer    int    `mapstructure:"bf2_mapsize"`
}

type StatusResponse struct {
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

func Status(ctx context.Context, conn net.Conn) (*StatusResponse, error) {
	deadline, ok := ctx.Deadline()
	if !ok {
		deadline = time.Now().Add(1 * time.Second)
	}

	err := conn.SetReadDeadline(deadline)
	if err != nil {
		return nil, err
	}

	_, err = conn.Write(query)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer([]byte{})
	b := make([]byte, 1024)
	for {
		read, err := conn.Read(b)
		if err != nil {
			if e, ok := err.(net.Error); !ok || !e.Timeout() {
				return nil, err
			}
			break
		}

		buf.Write(b[:read])
	}

	data := bytes.SplitN(buf.Bytes()[16:buf.Len()-2], delim, 2)[0]

	splitData := bytes.Split(data, []byte{0x00})

	out := make(map[string]string)

	for i := 0; i < len(splitData)-1; i += 2 {
		out[string(splitData[i])] = string(splitData[i+1])
	}

	var resp StatusResponse
	err = mapstructure.WeakDecode(out, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
