package prism

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMessage(t *testing.T) {
	tests := []struct {
		subject Subject
		fields  [][]byte
		output  []byte
	}{
		{
			subject: Subject("example-subject"),
			fields: [][]byte{
				[]byte("Hello, World!"),
			},
			output: []byte{0x1, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2d, 0x73, 0x75, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x2, 0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x2c, 0x20, 0x57, 0x6f, 0x72, 0x6c, 0x64, 0x21, 0x4, 0x0},
		},
		{
			subject: Subject("example-subject-2"),
			fields: [][]byte{
				{},
			},
			output: []byte{0x1, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2d, 0x73, 0x75, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x2d, 0x32, 0x2, 0x4, 0x0},
		},
	}

	for _, test := range tests {
		t.Run(string(test.subject), func(t *testing.T) {
			msg := NewMessage(test.subject, test.fields...)
			encoded := msg.Encode()
			require.Equal(t, test.output, encoded)

			decoded, err := DecodeMessage(encoded)
			require.NoError(t, err)

			require.Equal(t, test.subject, decoded.Subject)
			require.Equal(t, test.fields, decoded.Fields)
		})
	}
}
