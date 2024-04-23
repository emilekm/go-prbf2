package main

import (
	"context"
	"log"

	"github.com/emilekm/go-prbf2/prism"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

const (
	ipport   = "127.0.0.1:4712"
	username = "superuser"
	password = "******"
)

func run() error {
	ctx := context.Background()

	client, err := prism.Dial(ipport)
	if err != nil {
		return err
	}

	err = client.Auth.Login(ctx, username, password)
	if err != nil {
		return err
	}

	return nil
}
