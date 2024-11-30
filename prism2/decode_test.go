package prism2_test

import (
	"testing"

	"github.com/emilekm/go-prbf2/prism2"
	"github.com/stretchr/testify/require"
)

func TestUnmarshalMessage(t *testing.T) {
	tests := []struct {
		rawMsg []byte
		output any
		into   any
	}{
		{
			rawMsg: []byte("hash\x03serverchallenge"),
			output: &prism2.Login1Response{
				Hash:            []byte("hash"),
				ServerChallenge: []byte("serverchallenge"),
			},
			into: &prism2.Login1Response{},
		},
		{
			rawMsg: []byte("0\x031567934982\x03channel\x03playername\x03content\x0A1\x031567934982\x03channel\x03playername\x03content"),
			output: &prism2.ChatMessages{
				{
					Type:       prism2.ChatMessageTypeOpfor,
					Timestamp:  1567934982,
					Channel:    "channel",
					PlayerName: "playername",
					Content:    "content",
				},
				{
					Type:       prism2.ChatMessageTypeBlufor,
					Timestamp:  1567934982,
					Channel:    "channel",
					PlayerName: "playername",
					Content:    "content",
				},
			},
			into: &prism2.ChatMessages{},
		},
	}

	for _, test := range tests {
		err := prism2.Unmarshal(test.rawMsg, test.into)
		require.NoError(t, err)

		require.Equal(t, test.output, test.into)
	}
}
