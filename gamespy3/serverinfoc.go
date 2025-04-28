package gamespy3

import (
	"bytes"
	"context"
	"io"
	"time"
)

func (c *Client) ServerInfoC(ctx context.Context, fields []FieldType) ([]string, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	deadline, ok := ctx.Deadline()
	if !ok {
		deadline = time.Now().Add(1 * time.Second)
	}

	c.conn.SetReadDeadline(deadline)

	t := timestamp()

	fieldsQuery := []byte{}
	for _, f := range fields {
		fieldsQuery = append(fieldsQuery, byte(f))
	}

	query := bytes.Join([][]byte{
		{0xFE, 0xFD, 0x00},
		t,
		{byte(len(fieldsQuery))},
		fieldsQuery,
		{0x00, 0x00},
	}, []byte{})

	_, err := c.conn.Write(query)
	if err != nil {
		return nil, err
	}

	b := make([]byte, 1400)

	read, _, _, _, err := c.conn.ReadMsgUDP(b, nil)
	if err != nil {
		return nil, err
	}

	resp := make([]string, len(fields))

	const headerSize = 5
	buf := bytes.NewBuffer(b[headerSize:read])
	for i := range fields {
		str, err := buf.ReadString(0x00)
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			resp[i] = str
			break
		}
		resp[i] = str[:len(str)-1]
	}

	return resp, nil
}
