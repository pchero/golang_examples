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

	// Create a durable consumer
	_, err = js.AddConsumer("TEST_STREAM", &nats.ConsumerConfig{
		Durable:   "my-consumer",
		AckPolicy: nats.AckExplicitPolicy,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Subscribe to the subject
	_, err = js.QueueSubscribe("test.subject", "workers", func(msg *nats.Msg) {
		// Process the message
		log.Printf("Received message: %s\n", string(msg.Data))

		// Acknowledge the message
		msg.Ack()
	}, nats.Durable("my-consumer"))
	if err != nil {
		log.Fatal(err)
	}

	// Keep the connection alive
	select {}
}
