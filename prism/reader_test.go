package prism

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReader(t *testing.T) {
	buf := bytes.NewBuffer([]byte("\x01login1\x02test\x04\x00"))
	bufReader := bufio.NewReader(buf)

	r := NewReader(bufReader)
	rawMsg, err := r.ReadMessage()
	require.NoError(t, err)

	assert.Equal(t, Subject("login1"), rawMsg.Subject())
	assert.Equal(t, []byte("test"), rawMsg.Content())
}
