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

	serverAddr, err := net.ResolveUDPAddr("udp", ipport)
	if err != nil {
		return err
	}

	udp, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		return err
	}

	client := gamespy3.New(udp)

	status, err := client.ServerInfoB(ctx)
	if err != nil {
		return err
	}

	fmt.Printf("%+v", status)

	return nil
}
