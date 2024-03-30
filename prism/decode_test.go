package prism_test

import (
	"testing"

	"github.com/emilekm/go-prbf2/prism"
	"github.com/stretchr/testify/require"
)

func TestDecode(t *testing.T) {
	tests := []struct {
		payload []byte
		output  prism.Message
	}{
		{
			payload: []byte("\x01login1\x02hash\x03serverchallenge\x04\x00"),
			output: &prism.Login1Response{
				Hash:            []byte("hash"),
				ServerChallenge: []byte("serverchallenge"),
			},
		},
		{
			payload: []byte("\x01chat\x020\x031567934982\x03channel\x03playername\x03content\x0A1\x031567934982\x03channel\x03playername\x03content\x04\x00"),
			output: &prism.ChatMessages{
				{
					Type:       prism.ChatMessageTypeOpfor,
					Timestamp:  1567934982,
					Channel:    "channel",
					PlayerName: "playername",
					Content:    "content",
				},
				{
					Type:       prism.ChatMessageTypeBlufor,
					Timestamp:  1567934982,
					Channel:    "channel",
					PlayerName: "playername",
					Content:    "content",
				},
			},
		},
	}

	for _, test := range tests {
		msg, err := prism.Decode(test.payload)
		require.NoError(t, err)

		require.Equal(t, test.output, msg)
	}
}
