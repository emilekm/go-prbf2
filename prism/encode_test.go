package prism_test

import (
	"testing"

	"github.com/emilekm/go-prbf2/prism"
	"github.com/stretchr/testify/require"
)

type nestedMessage struct {
	TestField1 *int
}

type testMessage struct {
	TestField1 []byte
	TestField2 string
	TestField3 int
	TestField4 uint
	Nested     nestedMessage
}

func (m testMessage) Subject() prism.Subject {
	return "test"
}

func TestEncode(t *testing.T) {
	msg := testMessage{
		TestField1: []byte("sha"),
		TestField2: "hash",
		TestField3: -123,
		TestField4: 12,
		Nested: nestedMessage{
			TestField1: pointer(123),
		},
	}

	data, err := prism.Encode(msg)
	require.NoError(t, err)

	require.Equal(t, []byte("\x01test\x02sha\x03hash\x03-123\x0312\x03123\x04\x00"), data)
}

func pointer[T any](v T) *T {
	return &v
}
