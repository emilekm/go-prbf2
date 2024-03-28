package prism

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClient(t *testing.T) {
	host := os.Getenv("PRISM_HOST")
	port := os.Getenv("PRISM_PORT")
	user := os.Getenv("PRISM_USER")
	pass := os.Getenv("PRISM_PASS")

	if host == "" || port == "" || user == "" || pass == "" {
		t.Skip("skipping test; environment variables not set")
	}

	c := NewClient(ClientConfig{
		Host:     host,
		Port:     port,
		Username: user,
		Password: pass,
	})

	err := c.Connect()
	require.NoError(t, err)
}
