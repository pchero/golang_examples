package main

import (
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
)

func main() {
	// Connect to NATS server
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	// Subject to subscribe to
	subject := "request.subject"

	// Subscribe to the subject
	nc.Subscribe(subject, func(msg *nats.Msg) {
		// Process the request
		fmt.Printf("Received request: %s\n", string(msg.Data))

		// Respond to the request
		replyMsg := "Hi, this is a reply!"
		nc.Publish(msg.Reply, []byte(replyMsg))
	})

	// Keep the connection alive
	select {}
}
