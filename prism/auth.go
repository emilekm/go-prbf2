package prism

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"math/rand"
)

func auth(r *Responder, config ClientConfig) error {
	cck := cck(32)

	login1Req := Login1Request{
		ServerVersion:      ServerVersion1,
		Username:           config.Username,
		ClientChallengeKey: cck,
	}

	resp, err := r.Send(Marshal(login1Req), &SendOpts{
		ResponseSubjects: []Subject{SubjectLogin1},
	})
	if err != nil {
		return fmt.Errorf("login1: %w", err)
	}

	var login1Resp Login1Response
	err = UnmarshalInto(resp.Messages[0], &login1Resp)
	if err != nil {
		return fmt.Errorf("login1: %w", err)
	}

	passwordHash := sha1.New()
	saltedPassword := sha1.New()
	challengeDigest := sha1.New()

	_, err = passwordHash.Write([]byte(config.Password))
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
				[]byte(config.Username),
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

	_, err = r.Send(Marshal(login2Req), &SendOpts{
		ResponseSubjects: []Subject{SubjectConnected},
	})
	if err != nil {
		return fmt.Errorf("login2: %w", err)
	}

	return nil
}

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func cck(n int) []byte {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return b
}
