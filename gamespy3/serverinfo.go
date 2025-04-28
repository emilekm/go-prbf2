package gamespy3

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/hex"
	"io"
	"strings"
	"time"

	"github.com/go-viper/mapstructure/v2"
	"github.com/pkg/errors"
)

func (c *Client) ServerInfoB(ctx context.Context) (*ServerInformation, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	deadline, ok := ctx.Deadline()
	if !ok {
		deadline = time.Now().Add(1 * time.Second)
	}

	err := c.conn.SetReadDeadline(deadline)
	if err != nil {
		return nil, err
	}

	t := timestamp()

	query := bytes.Join([][]byte{
		{0xFE, 0xFD, 0x00},
		t,
		{0xFF, 0xFF, 0xFF},
		{0x01},
	}, []byte{})

	_, err = c.conn.Write(query)
	if err != nil {
		return nil, err
	}

	b := make([]byte, 1400)
	header := make(map[string]string)
	playersData := make(map[string][]string)
	teams := make(map[string][]string)
	hasMore := true
	for hasMore {
		read, _, _, _, err := c.conn.ReadMsgUDP(b, nil)
		if err != nil {
			return nil, err
		}

		println(hex.Dump(b[:read]))

		reader := bytes.NewReader(b[:read])

		_, err = reader.Seek(14, io.SeekCurrent)
		if err != nil {
			return nil, err
		}
		msgFlag, err := reader.ReadByte()
		if err != nil {
			return nil, err
		}

		if (msgFlag & 0x80) != 0 {
			hasMore = false
		}

		for reader.Len() > 0 {
			typ, err := reader.ReadByte()
			if err != nil {
				return nil, err
			}

			switch typ {
			case 0x00:
				err = unmarshalHeader(reader, header)
				if err != nil {
					return nil, errors.Wrap(err, "failed to unmarshal header")
				}
			case 0x01:
				err = unmarshalPlayers(reader, playersData)
				if err != nil {
					return nil, errors.Wrap(err, "failed to unmarshal players")
				}
			case 0x02:
				err = unmarshalTeam(reader, teams)
				if err != nil {
					return nil, errors.Wrap(err, "failed to unmarshal team")
				}
			}
		}
	}

	var h Header
	err = mapstructure.WeakDecode(header, &h)
	if err != nil {
		return nil, err
	}

	players := make([]Player, len(playersData["pid_"]))

	for k, v := range playersData {
		for i, p := range v {
			switch k {
			case "pid_":
				players[i].ProfileID = p
			case "skill_":
				players[i].Skill = p
			case "AIBot_":
				players[i].Bot = p == "1"
			case "team_":
				players[i].Team = p
			case "deaths_":
				players[i].Deaths = p
			case "score_":
				players[i].Score = p
			case "ping_":
				players[i].Ping = p
			case "player_":
				players[i].Player = p
			}
		}
	}

	return &ServerInformation{
		Header:  h,
		Players: players,
	}, nil
}

func unmarshalHeader(b *bytes.Reader, header map[string]string) error {
	for b.Len() > 0 {
		var values [2]strings.Builder
		for i := range values {
			for {
				c, err := b.ReadByte()
				if err != nil {
					return err
				}
				if c == 0x00 {
					break
				}
				err = values[i].WriteByte(c)
				if err != nil {
					return err
				}
			}
		}

		header[values[0].String()] = values[1].String()

		end := make([]byte, 2)
		read, err := b.Read(end)
		if err != nil {
			return err
		}

		if read == 2 && (end[1] == 0x01 || end[1] == 0x02) {
			b.UnreadByte()
			break
		}

		b.UnreadByte()
		b.UnreadByte()
	}

	return nil
}

func unmarshalPlayers(b *bytes.Reader, players map[string][]string) error {
	for b.Len() > 0 {

		var key strings.Builder

		for {
			c, err := b.ReadByte()
			if err != nil {
				return err
			}

			if c == 0x00 {
				// Advance by one more zero
				b.ReadByte()
				break
			}

			err = key.WriteByte(c)
			if err != nil {
				return err
			}
		}

		values := make([]strings.Builder, 0)
		var currentValue strings.Builder
		for {
			c, err := b.ReadByte()
			if err != nil {
				return err
			}

			if c == 0x00 {
				c_, err := b.ReadByte()
				if err != nil && err != io.EOF {
					return err
				}

				if c_ == 0x00 {
					break
				}

				err = b.UnreadByte()
				if err != nil {
					return err
				}

				values = append(values, currentValue)
				currentValue.Reset()
				continue
			}

			err = currentValue.WriteByte(c)
			if err != nil {
				return err
			}
		}

		if _, ok := players[key.String()]; !ok {
			players[key.String()] = make([]string, 0)
		}

		for _, v := range values {
			players[key.String()] = append(players[key.String()], v.String())
		}

		end := make([]byte, 2)
		read, err := b.Read(end)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		if read == 2 && end[1] == 0x02 {
			b.UnreadByte()
			break
		}

		b.UnreadByte()
		b.UnreadByte()
	}

	return nil
}

func unmarshalTeam(b *bytes.Reader, teams map[string][]string) error {
	buf := make([]byte, 512)
	read, err := b.Read(buf)
	if err != nil {
		return err
	}

	println(hex.Dump(buf[:read]))
	return nil
}

func timestamp() []byte {
	timestamp := uint32(time.Now().Unix())
	timestampBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(timestampBytes, timestamp)
	return timestampBytes
}
