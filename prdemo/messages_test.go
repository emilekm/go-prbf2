package prdemo

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testValRead struct {
	A uint8
}

func (t testValRead) Read(m *Message) (any, error) {
	err := m.Decode(&t.A)
	if err != nil {
		return t, err
	}
	return t, nil
}

func (t *testValRead) Decode(m *Message) error {
	tN, err := t.Read(m)
	if err != nil {
		return err
	}
	*(t) = tN.(testValRead)
	return nil
}

type testS struct {
	A uint8
	B *uint8
	C string
	D testValRead
	E *testValRead
}

func TestMessageDecode(t *testing.T) {
	s := testS{}

	buf := []byte{0x00, 0x02, 0x03}
	buf = append(buf, []byte("test")...)
	buf = append(buf, []byte{0x00, 0x02, 0x03}...)
	m, err := NewMessage(bytes.NewReader(buf))
	require.NoError(t, err)
	assert.Equal(t, ServerDetailsType, m.Type)
	err = m.Decode(&s)
	require.NoError(t, err)
	b := uint8(3)
	assert.Equal(t, testS{
		A: 0x02,
		B: &b,
		C: "test",
		D: testValRead{A: 2},
		E: &testValRead{A: 3},
	}, s)
	m, err = NewMessage(bytes.NewReader(buf))
	require.NoError(t, err)

	var p uint8
	pa := &s.A

	err = m.Decode(&p)
	require.NoError(t, err)
	require.Equal(t, uint8(2), p)
	err = m.Decode(pa)
	require.NoError(t, err)
	require.Equal(t, uint8(3), s.A)
}

const testMsgFmt = "testdata/messages/%s.bin"

func TestServerDetailsDecode(t *testing.T) {
	var details ServerDetails

	buf, err := os.ReadFile(fmt.Sprintf(testMsgFmt, "serverdetails"))
	require.NoError(t, err)
	msg, err := NewMessage(bytes.NewReader(buf))
	require.NoError(t, err)
	require.Equal(t, ServerDetailsType, msg.Type)

	err = msg.Decode(&details)
	require.NoError(t, err)

	sd := ServerDetails{Version: 3, DemoTimePerTick: 0.3, IPPort: "127.0.0.1:16567", ServerName: "[PR v1.7.4.5] Test Server", MaxPlayers: 0x64, RoundLength: 0x3840, BriefingTime: 0xf0, Map: Map{Name: "fields_of_kassel", Gamemode: "gpm_cq", Layer: 0x80}, BluforTeam: "ww2ger", OpforTeam: "ww2usa", StartTime: 0x62f3de20, Tickets1: 0x320, Tickets2: 0x2ee, MapSize: 4}
	require.Equal(t, sd, details)
}

func TestPlayerUpdatesDecode(t *testing.T) {
	var updates PlayersUpdate

	buf, err := os.ReadFile(fmt.Sprintf(testMsgFmt, "playerupdates"))
	require.NoError(t, err)
	msg, err := NewMessage(bytes.NewReader(buf))
	require.NoError(t, err)
	require.Equal(t, PlayerUpdateType, msg.Type)

	err = msg.Decode(&updates)
	require.NoError(t, err)

	// require.Equal(t, PlayersUpdate{}, updates)
}

func TestVehicleUpdatesDecode(t *testing.T) {
	var updates VehiclesUpdate

	buf, err := os.ReadFile(fmt.Sprintf(testMsgFmt, "vehicleupdates"))
	require.NoError(t, err)
	msg, err := NewMessage(bytes.NewReader(buf))
	require.NoError(t, err)
	require.Equal(t, VehicleUpdateType, msg.Type)

	err = msg.Decode(&updates)
	require.NoError(t, err)

	// require.Equal(t, VehiclesUpdate{}, updates)
}
