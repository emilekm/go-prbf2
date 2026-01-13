package logs

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestParsePlayerProfileEntry(t *testing.T) {
	tests := []struct {
		input  string
		output *PlayerProfileEntry
	}{
		{
			input: "[2025-01-05 19:53:02]	30f693875976497eafebe93691658449	2	cassius23",
			output: &PlayerProfileEntry{
				Timestamp:  time.Date(2025, 1, 5, 19, 53, 2, 0, time.UTC),
				KeyHash:    "30f693875976497eafebe93691658449",
				TrustLevel: 2,
				Username:   "cassius23",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			entry, err := ParsePlayerProfileEntry(test.input)
			require.NoError(t, err)
			require.Equal(t, test.output, entry)
		})
	}
}
