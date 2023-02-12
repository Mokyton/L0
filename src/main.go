package main

import (
	"github.com/nats-io/nats.go"
	"log"
	"time"
	//"time"
)

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	// Subscribe
	sub, err := nc.SubscribeSync("foo")
	if err != nil {
		log.Fatal(err)
	}

	// Wait for a message
	msg, err := sub.NextMsg(10 * time.Second)
	if err != nil {
		log.Fatal(err)
	}
	msg.Respond([]byte("Check"))

	log.Printf("Reply: %s", string(msg.Data))
}
