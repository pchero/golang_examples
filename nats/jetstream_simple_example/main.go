package main

import (
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)

	// Connect to NATS server
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		logrus.Errorf("Could not connect to the nats. err: %v", err)
		log.Fatal(err)
	}
	defer nc.Close()

	// Create JetStream context
	js, err := nc.JetStream()
	if err != nil {
		logrus.Errorf("Could not initiate the jetstream. err: %v", err)
		log.Fatal(err)
	}

	// Define stream configuration
	streamConfig := &nats.StreamConfig{
		Name:     "REQUESTS",
		Subjects: []string{"request.created"},
	}

	// Check if the stream already exists
	if _, err = js.StreamInfo("REQUESTS"); err == nil {
		// If the stream exists, update it
		_, err = js.UpdateStream(streamConfig)
		if err != nil {
			logrus.Errorf("Could not update the stream. err: %v", err)
			log.Fatal(err)
		}
	} else {
		// If the stream doesn't exist, create it
		_, err = js.AddStream(streamConfig)
		if err != nil {
			logrus.Errorf("Could not add the stream. err: %v", err)
			log.Fatal(err)
		}
	}

	logrus.Debugf("Adding consumer...")

	// Define consumer configuration for responders
	consumerConfig := &nats.ConsumerConfig{
		Durable:        "REQUESTS_CONSUMER",
		DeliverSubject: "deliver.requests",
		AckPolicy:      nats.AckExplicitPolicy,
	}

	// Add consumer
	_, err = js.AddConsumer("REQUESTS", consumerConfig)
	if err != nil {
		logrus.Errorf("Could not add the consumer. err: %v", err)
		log.Fatal(err)
	}

	// Start the responder
	go startResponder(js)

	// Give the responder some time to start
	time.Sleep(2 * time.Second)

	// Send a request
	sendRequest(js)

	// Wait to receive responses
	time.Sleep(5 * time.Second)
}

func startResponder(js nats.JetStreamContext) {
	// Subscribe to the consumer in a queue group
	sub, err := js.QueueSubscribe("deliver.requests", "responder-group", func(msg *nats.Msg) {
		fmt.Printf("Responder received request: %s\n", string(msg.Data))
		// Send a response
		response := fmt.Sprintf("response to '%s'", string(msg.Data))
		msg.Respond([]byte(response))
		// Acknowledge the message
		msg.Ack()
	}, nats.Durable("REQUESTS_CONSUMER"), nats.ManualAck())
	if err != nil {
		logrus.Errorf("Could not subscribe the queue. err: %v", err)

		log.Fatal(err)
	}
	defer sub.Unsubscribe()

	// Keep the responder running
	select {}
}

func sendRequest(js nats.JetStreamContext) {
	// Create an inbox for receiving the response
	inbox := nats.NewInbox()

	// Subscribe to the inbox
	sub, err := js.SubscribeSync(inbox)
	if err != nil {
		log.Fatal(err)
	}
	defer sub.Unsubscribe()

	// Send a request
	request := "request #1"
	msg := &nats.Msg{
		Subject: "request.created",
		Reply:   inbox,
		Data:    []byte(request),
	}
	js.PublishMsg(msg)

	// Wait for a response
	responseMsg, err := sub.NextMsg(5 * time.Second)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Requester received response: %s\n", string(responseMsg.Data))
}
