package prism2

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriter(t *testing.T) {
	buf := &bytes.Buffer{}
	bufWriter := bufio.NewWriter(buf)

	w := NewWriter(bufWriter)
	err := w.WriteMessage(NewMessage(SubjectLogin1, []byte("test")))

	assert.NoError(t, err)
	assert.Equal(t, buf.String(), "\x01login1\x02test\x04\x00")
}
