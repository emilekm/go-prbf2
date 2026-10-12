package logs

import (
	"errors"
	"strings"
	"time"
)

const (
	DefaultAdminEntryDateFormat = "2006_01_02_15_04_05"
)

type AdminEntry struct {
	Timestamp time.Time
	Action    string
	Issuer    string
	Target    string
	Details   string
}

func ParseAdminEntry(line string, dateFormat string) (*AdminEntry, error) {
	if !strings.HasPrefix(line, "[") {
		return nil, errors.New("invalid format: expected line to start with '['")
	}

	closeBracket := strings.Index(line, "]")
	if closeBracket == -1 {
		return nil, errors.New("invalid format: missing closing ']' for timestamp")
	}

	timestampStr := line[1:closeBracket]
	timestamp, err := time.Parse(dateFormat, timestampStr)
	if err != nil {
		return nil, err
	}

	rest, ok := strings.CutPrefix(line[closeBracket+1:], " ")
	if !ok {
		return nil, errors.New("invalid format: expected ' ' after timestamp")
	}

	performedByIdx := strings.Index(rest, " performed by '")
	if performedByIdx == -1 {
		return nil, errors.New("invalid format: missing ' performed by '")
	}

	action := strings.TrimSpace(rest[:performedByIdx])

	rest = rest[performedByIdx+15:]

	detailsIdx := strings.Index(rest, "': ")
	if detailsIdx == -1 {
		return nil, errors.New("invalid format: missing ': ' before details")
	}

	admin := rest[:detailsIdx]
	details := rest[detailsIdx+3:]

	var target string
	if onIdx := strings.Index(admin, "' on '"); onIdx != -1 {
		target = admin[onIdx+6:]
		admin = admin[:onIdx]
	}

	return &AdminEntry{
		Timestamp: timestamp,
		Action:    action,
		Issuer:    admin,
		Target:    target,
		Details:   details,
	}, nil
}
