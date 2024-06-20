package prism

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"math/rand"
)

func (c *Client) Login(ctx context.Context, username, password string) error {
	cck := cckGen(32)

	resp, err := c.Command(ctx, CommandLogin1, &Login1Request{
		ServerVersion:      ServerVersion1,
		Username:           username,
		ClientChallengeKey: cck,
	}, SubjectLogin1)

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
		return fmt.Errorf("login2: challengedigest: %w", err)
	}

	resp, err = c.Command(ctx, CommandLogin2, &Login2Request{
		ChallengeDigest: challengeDigestHash,
	}, SubjectConnected)
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
