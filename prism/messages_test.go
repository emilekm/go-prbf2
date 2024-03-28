package prism

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogin1Request(t *testing.T) {
	// Create a new request
	req := Login1Request{
		ServerVersion:      ServerVersion1,
		Username:           "test",
		ClientChallengeKey: []byte("test"),
	}

	msg := Marshal(req)
	assert.Equal(t, msg.Subject, SubjectLogin1)
	assert.Equal(t, []byte{0x31, 0x3, 0x74, 0x65, 0x73, 0x74, 0x3, 0x74, 0x65, 0x73, 0x74}, msg.Content)
}
