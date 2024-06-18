package prism

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
	w.WriteMessage(RawMessage{
		subject: SubjectLogin1,
		content: []byte("test"),
	})

	assert.Equal(t, buf.String(), "\x01login1\x02test\x04\x00")
}
