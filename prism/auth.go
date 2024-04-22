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
	receiver := NewReceiver(&a.c.Reader)
	defer receiver.Close()
	responder := NewResponder(receiver, &a.c.Writer)

	cck := cckGen(32)

	login1Req := Login1Request{
		ServerVersion:      ServerVersion1,
		Username:           username,
		ClientChallengeKey: cck,
	}

	resp, err := responder.SendWithResponse(
		ctx,
		login1Req,
		ResponseWithMessageSubject(SubjectLogin1),
	)
	if err != nil {
		return fmt.Errorf("login1: %w", err)
	}

	var login1Resp Login1Response
	err = DecodeContent(resp.Messages[0].(RawMessage).Content(), &login1Resp)
	if err != nil {
		return fmt.Errorf("login1: %w", err)
	}

	passwordHash := sha1.New()
	saltedPassword := sha1.New()
	challengeDigest := sha1.New()

	_, err = passwordHash.Write([]byte(password))
	if err != nil {
		return err
	}

	salted := append(login1Resp.Hash, SeparatorStart...)
	salted = append(salted, hex.EncodeToString(passwordHash.Sum(nil))...)

	_, err = saltedPassword.Write(salted)
	if err != nil {
		return err
	}

	_, err = challengeDigest.Write(
		bytes.Join(
			[][]byte{
				[]byte(username),
				cck,
				login1Resp.ServerChallenge,
				[]byte(hex.EncodeToString(saltedPassword.Sum(nil))),
			},
			SeparatorField,
		),
	)
	if err != nil {
		return err
	}

	login2Req := Login2Request{
		ChallengeDigest: hex.EncodeToString(challengeDigest.Sum(nil)),
	}

	_, err = responder.SendWithResponse(
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
