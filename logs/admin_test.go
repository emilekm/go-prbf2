package logs

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestParseAdminEntry(t *testing.T) {
	tests := []struct {
		input  string
		output *AdminEntry
	}{
		{
			input: "[2026_01_09_23_25_47] MAPVOTERESULT   performed by 'TAG user': Vote finished: Ramiel: 14 | Dragon Fly: 3 | Sbeneh Outskirts: 15",
			output: &AdminEntry{
				Timestamp: time.Date(2026, 1, 9, 23, 25, 47, 0, time.UTC),
				Action:    "MAPVOTERESULT",
				Issuer:    "TAG user",
				Target:    "",
				Details:   "Vote finished: Ramiel: 14 | Dragon Fly: 3 | Sbeneh Outskirts: 15",
			},
		},
		{
			input: "[2026_01_09_23_25_56] !SETNEXT        performed by ' user': Sbeneh Outskirts (Insurgency, Std)",
			output: &AdminEntry{
				Timestamp: time.Date(2026, 1, 9, 23, 25, 56, 0, time.UTC),
				Action:    "!SETNEXT",
				Issuer:    " user",
				Target:    "",
				Details:   "Sbeneh Outskirts (Insurgency, Std)",
			},
		},
		{
			input: "[2026_01_09_23_26_40] !REPORTP        performed by 'TAG user1' on ' user2': afk",
			output: &AdminEntry{
				Timestamp: time.Date(2026, 1, 9, 23, 26, 40, 0, time.UTC),
				Action:    "!REPORTP",
				Issuer:    "TAG user1",
				Target:    " user2",
				Details:   "afk",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			entry, err := ParseAdminEntry(test.input, DefaultAdminEntryDateFormat)
			require.NoError(t, err)
			require.Equal(t, test.output, entry)
		})
	}
}
