package logs

import (
	"errors"
	"net"
	"strconv"
	"strings"
	"time"
)

type AccountStatus int

const (
	StatusEmpty AccountStatus = iota
	StatusLegacy
	StatusWhitelisted
	StatusVacBanned
)

type JoinEntry struct {
	Timestamp  time.Time
	KeyHash    string
	TrustLevel int
	Name       string
	CreatedAt  time.Time
	IP         net.IP
	Status     AccountStatus
}

func ParseJoinEntry(line string) (*JoinEntry, error) {
	parts := strings.Split(line, "\t")
	if len(parts) < 6 {
		return nil, errors.New("invalid format: expected at least 6 fields")
	}

	timestamp, err := time.Parse("[2006-01-02 15:04:05]", parts[0])
	if err != nil {
		return nil, err
	}

	trustLevel, err := strconv.Atoi(parts[2])
	if err != nil {
		return nil, err
	}

	createdAt, err := time.Parse("2006-01-02", parts[4])
	if err != nil {
		return nil, err
	}

	ip := net.ParseIP(parts[5])
	if ip == nil {
		return nil, errors.New("invalid IP address")
	}

	status := StatusEmpty
	if len(parts) > 6 {
		switch parts[6] {
		case "(LEGACY)":
			status = StatusLegacy
		case "(WHITELISTED)":
			status = StatusWhitelisted
		case "(VAC BANNED)":
			status = StatusVacBanned
		}
	}

	return &JoinEntry{
		Timestamp:  timestamp,
		KeyHash:    parts[1],
		TrustLevel: trustLevel,
		Name:       parts[3],
		CreatedAt:  createdAt,
		IP:         ip,
		Status:     status,
	}, nil
}
