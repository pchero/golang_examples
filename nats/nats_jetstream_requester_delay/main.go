package main

import (
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

	// Connect to JetStream context
	js, err := nc.JetStream()
	if err != nil {
		log.Fatal(err)
	}

	// Define the subject
	subject := "request.subject"

	// Message to be sent
	message := []byte("This is a delayed message")

	// Delay duration
	delay := 10 * time.Second

	// Publish the message with a delay
	msg := &nats.Msg{
		Subject: subject,
		Data:    message,
		Header:  nats.Header{"Nats-Delay": {delay.String()}},
	}
	_, err = js.PublishMsg(msg)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Published delayed message to '%s', delay: %s\n", subject, delay)
}
