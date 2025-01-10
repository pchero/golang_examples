package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/nats-io/nats.go"
)

func main() {
	// Connect to NATS server
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	// Queue group name
	queueName := "workers"

	// Subscriber 1
	_, err = nc.QueueSubscribe("request.subject", queueName, func(msg *nats.Msg) {
		log.Printf("[Subscriber 1] Received message: %s\n", string(msg.Data))
		// Respond to the request
		replyMsg := "Hi, this is a reply!"
		nc.Publish(msg.Reply, []byte(replyMsg))
	})
	if err != nil {
		log.Fatal(err)
	}

	// // Subscriber 2
	// _, err = nc.QueueSubscribe("request.subject", queueName, func(msg *nats.Msg) {
	// 	log.Printf("[Subscriber 2] Received message: %s\n", string(msg.Data))

	// 	// Respond to the request
	// 	replyMsg := "Hi, this is a reply!"
	// 	nc.Publish(msg.Reply, []byte(replyMsg))

	// })
	// if err != nil {
	// 	log.Fatal(err)
	// }

	log.Println("Subscribers are ready to receive messages...")

	// Wait for interrupt signal to gracefully shut down
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
}
