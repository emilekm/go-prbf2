package bf2demo

import (
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

type Metadata struct {
	ServerName string
	StartTime  string
	MapName    string
}

var demoEndian = binary.LittleEndian

func Open(path string) (io.ReadCloser, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	return zlib.NewReader(file)
}

func DecodeMetadata(r io.Reader) (*Metadata, error) {
	var metadata Metadata

	err := binary.Read(r, demoEndian, make([]byte, 4)) // skip first 4 bytes
	if err != nil {
		return nil, fmt.Errorf("failed to read first 4 bytes: %w", err)
	}

	metadata.ServerName, err = readString32(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read server name: %w", err)
	}

	metadata.StartTime, err = readString32(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read start time: %w", err)
	}

	metadata.MapName, err = readString32(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read map name: %w", err)
	}

	return &metadata, nil
}

func readString32(r io.Reader) (string, error) {
	var length uint32
	err := binary.Read(r, demoEndian, &length)
	if err != nil {
		return "", fmt.Errorf("failed to read string length: %w", err)
	}

	buf := make([]byte, length)
	_, err = r.Read(buf)
	if err != nil {
		return "", fmt.Errorf("failed to read string: %w", err)
	}

	return string(buf), nil
}
