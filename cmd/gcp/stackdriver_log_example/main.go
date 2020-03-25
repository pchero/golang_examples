package main

import (
	"context"
	"flag"
	"log"

	"cloud.google.com/go/logging"
)

var flagProjectID = flag.String("project-id", "test-project", "project id")
var flagLogName = flag.String("log-name", "test-log", "log name")

func main() {
	ctx := context.Background()
	flag.Parse()

	// create a client
	client, err := logging.NewClient(ctx, *flagProjectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	logger := client.Logger(*flagLogName).StandardLogger(logging.Info)

	logger.Println("hello world")
}