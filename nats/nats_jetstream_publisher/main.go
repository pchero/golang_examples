package main

import (
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

	// Connect to JetStream context
	js, err := nc.JetStream()
	if err != nil {
		log.Fatal(err)
	}

	// Publish a message with a unique ID
	messageID := "unique-message-id"
	msg := nats.NewMsg("test.subject")
	msg.Header.Set("Nats-Msg-Id", messageID)
	msg.Data = []byte("Hello, JetStream!")

	_, err = js.PublishMsg(msg)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Published message with ID: %s\n", messageID)
}
