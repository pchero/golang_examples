package main

import (
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	// Connect to NATS server
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	// Subject to send request to
	subject := "request.subject"

	// Message to send
	requestMsg := "Hello, can you reply?"

	// Send request and wait for reply
	msg, err := nc.Request(subject, []byte(requestMsg), 2*time.Second)
	if err != nil {
		log.Fatal(err)
	}

	// Print the response
	fmt.Printf("Received reply: %s\n", string(msg.Data))
}
