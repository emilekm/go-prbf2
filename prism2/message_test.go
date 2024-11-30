package prism2_test

import (
	"testing"

	"github.com/emilekm/go-prbf2/prism2"
	"github.com/stretchr/testify/require"
)

func TestMessageUnmarshalBinary(t *testing.T) {
	tests := []struct {
		name   string
		buffer []byte
		msg    *prism2.Message
		err    bool
	}{
		{
			name:   "success",
			buffer: []byte("\x01login1\x02test\x04"),
			msg: prism2.NewMessage(
				prism2.SubjectLogin1,
				[]byte("test"),
			),
		},
		{
			name:   "invalid-message",
			buffer: []byte("invalid-message"),
			err:    true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			msg := &prism2.Message{}
			err := msg.UnmarshalBinary(test.buffer)
			if test.err {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, test.msg, msg)
			}
		})
	}
}

func TestMessageMarshalBinary(t *testing.T) {
	tests := []struct {
		name   string
		msg    *prism2.Message
		buffer []byte
	}{
		{
			name: "success",
			msg: prism2.NewMessage(
				prism2.SubjectLogin1,
				[]byte("test"),
			),
			buffer: []byte("\x01login1\x02test\x04"),
		},
		{
			name: "empty-body",
			msg: prism2.NewMessage(
				prism2.SubjectLogin1,
				[]byte{},
			),
			buffer: []byte("\x01login1\x02\x04"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			buffer, err := test.msg.MarshalBinary()
			require.NoError(t, err)
			require.Equal(t, test.buffer, buffer)
		})
	}
}
