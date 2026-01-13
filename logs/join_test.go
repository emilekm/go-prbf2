package logs

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestParseJoinEntry(t *testing.T) {
	tests := []struct {
		input  string
		output *JoinEntry
	}{
		{
			input: "[2026-01-09 23:02:08]	30f693875976497eafebe93691658449	0	TAG Username	2020-05-01	127.0.0.1",
			output: &JoinEntry{
				Timestamp:  time.Date(2026, 1, 9, 23, 2, 8, 0, time.UTC),
				KeyHash:    "30f693875976497eafebe93691658449",
				TrustLevel: 0,
				Name:       "TAG Username",
				CreatedAt:  time.Date(2020, 5, 1, 0, 0, 0, 0, time.UTC),
				IP:         net.ParseIP("127.0.0.1"),
			},
		},
		{
			input: "[2026-01-09 23:02:08]	30f693875976497eafebe93691658449	0	TAG Username	2020-05-01	127.0.0.1	(LEGACY)",
			output: &JoinEntry{
				Timestamp:  time.Date(2026, 1, 9, 23, 2, 8, 0, time.UTC),
				KeyHash:    "30f693875976497eafebe93691658449",
				TrustLevel: 0,
				Name:       "TAG Username",
				CreatedAt:  time.Date(2020, 5, 1, 0, 0, 0, 0, time.UTC),
				IP:         net.ParseIP("127.0.0.1"),
				Status:     StatusLegacy,
			},
		},
		{
			input: "[2026-01-09 23:02:08]	30f693875976497eafebe93691658449	1	TAG Username	2020-05-01	127.0.0.1	(WHITELISTED)",
			output: &JoinEntry{
				Timestamp:  time.Date(2026, 1, 9, 23, 2, 8, 0, time.UTC),
				KeyHash:    "30f693875976497eafebe93691658449",
				TrustLevel: 1,
				Name:       "TAG Username",
				CreatedAt:  time.Date(2020, 5, 1, 0, 0, 0, 0, time.UTC),
				IP:         net.ParseIP("127.0.0.1"),
				Status:     StatusWhitelisted,
			},
		},
		{
			input: "[2026-01-09 23:02:08]	30f693875976497eafebe93691658449	2	 Username	2020-05-01	127.0.0.1	(VAC BANNED)",
			output: &JoinEntry{
				Timestamp:  time.Date(2026, 1, 9, 23, 2, 8, 0, time.UTC),
				KeyHash:    "30f693875976497eafebe93691658449",
				TrustLevel: 2,
				Name:       " Username",
				CreatedAt:  time.Date(2020, 5, 1, 0, 0, 0, 0, time.UTC),
				IP:         net.ParseIP("127.0.0.1"),
				Status:     StatusVacBanned,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result, err := ParseJoinEntry(test.input)
			require.NoError(t, err)
			require.Equal(t, test.output, result)
		})
	}
}
