package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/nats-io/nats.go"
)

// In-memory store to track processed messages
var processedMessages = make(map[string]struct{})
var mu sync.Mutex

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
		// Use the message ID (or create one) as the deduplication key
		messageID := string(msg.Data) // Simplified for example; use a proper ID in real applications

		// Check if the message has already been processed
		mu.Lock()
		if _, exists := processedMessages[messageID]; exists {
			mu.Unlock()
			fmt.Printf("Duplicate message received and ignored: %s\n", messageID)
			return
		}
		// Mark the message as processed
		processedMessages[messageID] = struct{}{}
		mu.Unlock()

		// Process the message
		fmt.Printf("Processing message: %s\n", messageID)

		// (Optional) Clear old entries to prevent memory bloat
	})

	// Keep the connection alive
	select {}
}
