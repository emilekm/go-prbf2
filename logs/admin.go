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

	rest := line[closeBracket+2:]

	performedByIdx := strings.Index(rest, " performed by '")
	if performedByIdx == -1 {
		return nil, errors.New("invalid format: missing ' performed by '")
	}

	action := strings.TrimSpace(rest[:performedByIdx])

	rest = rest[performedByIdx+15:]

	adminEnd := strings.Index(rest, "'")
	if adminEnd == -1 {
		return nil, errors.New("invalid format: missing closing quote for admin")
	}

	admin := rest[:adminEnd]
	rest = rest[adminEnd+1:]

	var target, details string

	if strings.HasPrefix(rest, " on '") {
		rest = rest[5:]
		targetEnd := strings.Index(rest, "'")
		if targetEnd == -1 {
			return nil, errors.New("invalid format: missing closing quote for target")
		}
		target = rest[:targetEnd]
		rest = rest[targetEnd+1:]
	}

	if strings.HasPrefix(rest, ": ") {
		details = rest[2:]
	} else {
		return nil, errors.New("invalid format: missing ': ' before details")
	}

	return &AdminEntry{
		Timestamp: timestamp,
		Action:    action,
		Issuer:    admin,
		Target:    target,
		Details:   details,
	}, nil
}
