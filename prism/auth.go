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

	login1Msg := NewMessage(
		SubjectLogin1,
		[]byte("1"),
		[]byte(config.User),
		cck,
	)

	resp, err := r.Send(login1Msg, &SendOpts{
		ResponseSubjects: []Subject{SubjectLogin1},
	})
	if err != nil {
		return fmt.Errorf("login1: %w", err)
	}

	passHash := resp.Messages[0].Fields[0]
	serverChallenge := resp.Messages[0].Fields[1]

	passwordHash := sha1.New()
	saltedPassword := sha1.New()
	challengeDigest := sha1.New()

	_, err = passwordHash.Write([]byte(config.Pass))
	if err != nil {
		return err
	}

	salted := append(passHash, SeparatorStart...)
	salted = append(salted, hex.EncodeToString(passwordHash.Sum(nil))...)

	_, err = saltedPassword.Write(salted)
	if err != nil {
		return err
	}

	_, err = challengeDigest.Write(
		bytes.Join(
			[][]byte{
				[]byte(config.User),
				cck,
				serverChallenge,
				[]byte(hex.EncodeToString(saltedPassword.Sum(nil))),
			},
			SeparatorField,
		),
	)
	if err != nil {
		return err
	}

	login2Msg := Message{
		Subject: SubjectLogin2,
		Fields:  [][]byte{[]byte(hex.EncodeToString(challengeDigest.Sum(nil)))},
	}

	_, err = r.Send(login2Msg, &SendOpts{
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
