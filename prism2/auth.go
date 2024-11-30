package prism2

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"math/rand"
)

const (
	SubjectLogin1    Subject = "login1"
	SubjectConnected Subject = "connected"

	CommandLogin1 Command = "login1"
	CommandLogin2 Command = "login2"
)

type Login1Request struct {
	ServerVersion      ServerVersion
	Username           string
	ClientChallengeKey []byte
}

type Login1Response struct {
	Hash            []byte
	ServerChallenge []byte
}

type Login2Request struct {
	ChallengeDigest string
}

// Login authenticates the client with given username and password.
// It is obligatory to call this method before any other method.
func (c *Client) Login(ctx context.Context, username, password string) error {
	cck := cckGen(32)

	login1Response, err := c.login1(ctx, username, cck)
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
		return fmt.Errorf("login2: challengedigest: %w", err)
	}

	return c.login2(ctx, challengeDigestHash)
}

func (c *Client) login1(ctx context.Context, username string, cck []byte) (*Login1Response, error) {
	login1Req, err := Marshal(&Login1Request{
		ServerVersion:      ServerVersion1,
		Username:           username,
		ClientChallengeKey: cck,
	})
	if err != nil {
		return nil, fmt.Errorf("login1: %w", err)
	}

	respCh := c.Send(&Request{
		Message:         NewMessage(CommandLogin1, login1Req),
		ExpectedSubject: SubjectLogin1,
	})

	for {
		select {
		case resp := <-respCh:
			if resp.Error != nil {
				return nil, fmt.Errorf("login1: %w", resp.Error)
			}

			var login1Response Login1Response
			err = Unmarshal(resp.Message.Body(), &login1Response)
			if err != nil {
				return nil, fmt.Errorf("login1: %w", err)
			}

			return &login1Response, nil
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

func (c *Client) login2(ctx context.Context, challengeDigestHash string) error {
	login2Req, err := Marshal(&Login2Request{
		ChallengeDigest: challengeDigestHash,
	})
	if err != nil {
		return fmt.Errorf("login2: %w", err)
	}

	respCh := c.Send(&Request{
		Message:         NewMessage(CommandLogin1, login2Req),
		ExpectedSubject: SubjectLogin1,
	})

	for {
		select {
		case resp := <-respCh:
			if resp.Error != nil {
				return fmt.Errorf("login2: %w", resp.Error)
			}

			var login1Response Login1Response
			err = Unmarshal(resp.Message.Body(), &login1Response)
			if err != nil {
				return fmt.Errorf("login2: %w", err)
			}

			return nil
		case <-ctx.Done():
			return ctx.Err()
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
