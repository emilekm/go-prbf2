package main

import (
	"context"
	"fmt"
	"log"
	"time"

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

	err = client.Login(ctx, username, password)
	if err != nil {
		return err
	}

	println("Logged in")

	details, err := client.ServerDetails(ctx)
	if err != nil {
		return err
	}

	fmt.Printf("Server details: %+v\n", details)

	players, err := client.ListPlayers(ctx)
	if err != nil {
		return err
	}

	for _, player := range players {
		fmt.Printf("Player: %+v\n", player)
	}

	sub := client.SubscribeAll()

	for msg := range sub {
		fmt.Printf("Message: %+v\n", msg.Subject())

		switch msg.Subject() {
		case prism.SubjectUpdatePlayers:
			var players prism.UpdatePlayers
			err := prism.Unmarshal(msg.Body(), &players)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
