package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	// Connect to NATS server
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	// Connect to JetStream context
	js, err := nc.JetStream()
	if err != nil {
		logrus.Errorf("Could not connect to the stream. err: %v", err)
		return
	}

	initJetstream(js)
	initConsumer(js)

	// Queue group name
	queueName := "workers"

	// // Subscribe to the subject
	// _, err = js.QueueSubscribe("updates", queueName, func(msg *nats.Msg) {
	// 	log.Printf("Received message: %s\n", string(msg.Data))
	// })
	// if err != nil {
	// 	logrus.Errorf("Could not subscribe the queue. err: %v", err)
	// 	return
	// }

	// Create durable consumer with queue group
	_, err = js.QueueSubscribe("updates", queueName, func(msg *nats.Msg) {
		log.Printf("Received message: %s\n", string(msg.Data))
		msg.Ack()
	}, nats.Durable("durable-worker"), nats.ManualAck())
	if err != nil {
		log.Fatalf("Could not subscribe to queue. err: %v", err)
	}

	log.Println("Subscriber is ready to receive messages...")

	// Wait for interrupt signal to gracefully shut down
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
}

func initJetstream(js nats.JetStreamContext) {
	// Check if the stream exists
	streamInfo, err := js.StreamInfo("TESTUPDATES")
	if err == nil {
		log.Printf("Stream %s already exists", streamInfo.Config.Name)
		return
	} else if err != nats.ErrStreamNotFound {
		log.Fatalf("Could not get nats jetstream info. err: %v", err)
	}

	if err == nats.ErrStreamNotFound {
		// Define the stream configuration
		streamConfig := &nats.StreamConfig{
			Name:     "TESTUPDATES",
			Subjects: []string{"updates"},
			Storage:  nats.FileStorage,
		}
		// Add the stream
		_, err = js.AddStream(streamConfig)
		if err != nil {
			log.Fatalf("Could not add the stream. err: %v", err)
		}
	} else {
		log.Fatalf("Error retrieving stream info: %v", err)
	}

	// // Define the stream configuration
	// streamConfig := &nats.StreamConfig{
	// 	Name:     "TESTUPDATES",
	// 	Subjects: []string{"updates"},
	// 	Storage:  nats.FileStorage,
	// }

	// // Add the stream
	// _, err = js.AddStream(streamConfig)
	// if err != nil {
	// 	logrus.Errorf("Could not add the stream. err: %v", err)
	// 	return
	// }

	// Define the consumer configuration
	consumerConfig := &nats.ConsumerConfig{
		Durable:   "durable-worker",
		AckPolicy: nats.AckExplicitPolicy,
	}

	// Add the consumer
	_, err = js.AddConsumer("TESTUPDATES", consumerConfig)
	if err != nil {
		logrus.Errorf("Could not add the consumer. err: %v", err)
		return
	}

}

func initConsumer(js nats.JetStreamContext) {
	// Define the consumer configuration
	consumerConfig := &nats.ConsumerConfig{
		Durable:        "durable-worker",
		AckPolicy:      nats.AckExplicitPolicy,
		DeliverSubject: "updates-push",
	}

	// Add the consumer
	_, err := js.AddConsumer("TESTUPDATES", consumerConfig)
	if err != nil && err != nats.ErrConsumerNameAlreadyInUse {
		log.Fatalf("Could not add the consumer. err: %v", err)
	}

}
