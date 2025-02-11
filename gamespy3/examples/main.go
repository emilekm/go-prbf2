package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/emilekm/go-prbf2/gamespy3"
)

const (
	ipport = "localhost:29900"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx := context.Background()

	conn, err := net.Dial("udp", ipport)
	if err != nil {
		return err
	}

	status, err := gamespy3.Status(ctx, conn)
	if err != nil {
		return err
	}

	fmt.Printf("%+v", status)

	return nil
}
