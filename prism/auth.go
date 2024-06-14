package prism

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"math/rand"
	"slices"
)

func (c *Client) Login(ctx context.Context, username, password string) error {
	cck := cckGen(32)

	login1Content := bytes.Join([][]byte{
		[]byte(ServerVersion1),
		[]byte(username),
		[]byte(cck),
	}, SeparatorField)

	id := c.Next()
	c.StartRequest(id)
	c.StartResponse(id)

	err := c.WriteMessage(&RawMessage{
		subject: SubjectLogin1,
		content: login1Content,
	})
	c.EndRequest(id)
	if err != nil {
		c.EndResponse(id)
		return fmt.Errorf("login1: %w", err)
	}

	login1Resp, err := c.waitForMessage(ctx, SubjectLogin1)
	if err != nil {
		c.EndResponse(id)
		return fmt.Errorf("login1: %w", err)
	}

	c.EndResponse(id)

	login1RespParts := bytes.SplitN(login1Resp.Content(), SeparatorField, 2)
	hash := login1RespParts[0]
	serverChallenge := login1RespParts[1]

	challengeDigestHash, err := prepareChallengeDigest(
		username,
		password,
		hash,
		cck,
		serverChallenge,
	)
	if err != nil {
		return fmt.Errorf("login1: challengedigest: %w", err)
	}

	id = c.Next()
	c.StartRequest(id)
	c.StartResponse(id)

	err = c.WriteMessage(&RawMessage{
		subject: SubjectLogin2,
		content: []byte(challengeDigestHash),
	})
	c.EndRequest(id)
	if err != nil {
		c.EndResponse(id)
		return fmt.Errorf("login2: %w", err)
	}

	_, err = c.waitForMessage(ctx, SubjectConnected)
	c.EndResponse(id)
	if err != nil {
		return fmt.Errorf("login2: %w", err)
	}

	return nil
}

func (c *Client) waitForMessage(ctx context.Context, expected Subject) (*RawMessage, error) {
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			msg, err := c.ReadMessage()
			if err != nil {
				return nil, err
			}

			if slices.Contains(errorSubjects, msg.Subject()) {
				return nil, ErrorMessageToError(msg)
			}

			if msg.Subject() == expected {
				return msg, nil
			}
		}
	}
}

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func cckGen(n int) []byte {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return b
}

func stringHash(s string) string {
	h := sha1.New()
	_, _ = h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func bytesHash(b ...[]byte) string {
	h := sha1.New()
	_, _ = h.Write(bytes.Join(b, nil))
	return hex.EncodeToString(h.Sum(nil))
}

func prepareChallengeDigest(username, password string, salt, clientChallenge, serverChallenge []byte) (string, error) {
	passwordHash := stringHash(password)

	saltedPasswordHash := bytesHash(salt, SeparatorStart, []byte(passwordHash))

	challengeDigestHash := bytesHash(
		bytes.Join([][]byte{
			[]byte(username),
			clientChallenge,
			serverChallenge,
			[]byte(saltedPasswordHash),
		}, SeparatorField),
	)

	return challengeDigestHash, nil
}
