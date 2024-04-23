package prism_test

import (
	"testing"

	"github.com/emilekm/go-prbf2/prism"
	"github.com/stretchr/testify/require"
)

func TestDecode(t *testing.T) {
	tests := []struct {
		rawMsg *prism.RawMessage
		output prism.Message
		into   prism.Message
	}{
		{
			rawMsg: prism.NewRawMessage("login1", []byte("hash\x03serverchallenge")),
			output: &prism.Login1Response{
				Hash:            []byte("hash"),
				ServerChallenge: []byte("serverchallenge"),
			},
			into: &prism.Login1Response{},
		},
		{
			rawMsg: prism.NewRawMessage("chat", []byte("0\x031567934982\x03channel\x03playername\x03content\x0A1\x031567934982\x03channel\x03playername\x03content")),
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
			into: &prism.ChatMessages{},
		},
	}

	for _, test := range tests {
		err := prism.DecodeRawMessage(test.rawMsg, test.into)
		require.NoError(t, err)

		require.Equal(t, test.output, test.into)
	}
}
