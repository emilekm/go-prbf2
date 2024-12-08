package prism2

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
	msg, err := r.ReadMessage()
	require.NoError(t, err)

	assert.Equal(t, Subject("login1"), msg.subject)
	assert.Equal(t, []byte("test"), msg.body)
}
