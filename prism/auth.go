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

	login1Request := Login1Request{
		ServerVersion:      ServerVersion1,
		Username:           username,
		ClientChallengeKey: cck,
	}

	id := c.Next()
	c.StartRequest(id)
	c.StartResponse(id)

	err := c.WriteMessage(&login1Request)
	c.EndRequest(id)
	if err != nil {
		c.EndResponse(id)
		return fmt.Errorf("login1: %w", err)
	}

	resp, err := c.waitForMessage(ctx, SubjectLogin1)
	if err != nil {
		c.EndResponse(id)
		return fmt.Errorf("login1: %w", err)
	}

	c.EndResponse(id)

	var login1Response Login1Response
	err = UnmarshalMessage(resp.Body(), &login1Response)
	if err != nil {
		return fmt.Errorf("login1: %w", err)
	}

	challengeDigestHash, err := prepareChallengeDigest(
		username,
		password,
		login1Response.Hash,
		cck,
		login1Response.ServerChallenge,
	)
	if err != nil {
		return fmt.Errorf("login1: challengedigest: %w", err)
	}

	id = c.Next()
	c.StartRequest(id)
	c.StartResponse(id)

	err = c.WriteMessage(&Login2Request{
		ChallengeDigest: challengeDigestHash,
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
				var errMsg Error
				err := UnmarshalMessage(msg.Body(), &errMsg)
				if err != nil {
					return nil, err
				}
				return nil, errMsg
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
