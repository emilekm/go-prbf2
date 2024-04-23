package prism

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"math/rand"
)

type Auth struct {
	c *Client
}

func NewAuth(c *Client) *Auth {
	return &Auth{c: c}
}

func (a *Auth) Login(ctx context.Context, username, password string) error {
	cck := cckGen(32)

	login1Req := Login1Request{
		ServerVersion:      ServerVersion1,
		Username:           username,
		ClientChallengeKey: cck,
	}

	resp, err := a.c.Send(
		ctx,
		login1Req,
		ResponseWithMessageSubject(SubjectLogin1),
	)
	if err != nil {
		return fmt.Errorf("login1: %w", err)
	}

	var login1Resp Login1Response
	err = DecodeRawMessage(resp.Messages[0], &login1Resp)
	if err != nil {
		return fmt.Errorf("login1: %w", err)
	}

	challengeDigestHash, err := prepareChallengeDigest(
		username,
		password,
		login1Resp.Hash,
		cck,
		login1Resp.ServerChallenge,
	)
	if err != nil {
		return fmt.Errorf("login1: challengedigest: %w", err)
	}

	login2Req := Login2Request{
		ChallengeDigest: challengeDigestHash,
	}

	_, err = a.c.Send(
		ctx,
		login2Req,
		ResponseWithMessageSubject(SubjectConnected),
	)
	if err != nil {
		return fmt.Errorf("login2: %w", err)
	}

	return nil
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
