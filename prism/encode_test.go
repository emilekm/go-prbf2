package prism_test

import (
	"testing"

	"github.com/emilekm/go-prbf2/prism"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type nestedMessage struct {
	TestField1 *int
}

type testMessage struct {
	TestField1    []byte
	TestField2    string
	TestField3    int
	TestField4    uint
	TestFieldBool bool
	Nested        nestedMessage
}

func TestMarshalMessage(t *testing.T) {
	msg := testMessage{
		TestField1:    []byte("sha"),
		TestField2:    "hash",
		TestField3:    -123,
		TestField4:    12,
		TestFieldBool: false,
		Nested: nestedMessage{
			TestField1: pointer(123),
		},
	}

	body, err := prism.Marshal(&msg)
	require.NoError(t, err)

	assert.Equal(t, []byte("sha\x03hash\x03-123\x0312\x030\x03123"), body)
}

func pointer[T any](v T) *T {
	return &v
}
