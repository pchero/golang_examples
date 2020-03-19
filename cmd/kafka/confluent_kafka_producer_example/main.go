package main

import (
	"fmt"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
	fmt.Printf("Hello world\n")

	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <broker> <topic>\n", os.Args[0])
		os.Exit(1)
	}

	broker := os.Args[1]
	topic := os.Args[2]

	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.server": broker})
	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("Created producer. %v\n", p)

	// optional delivery channel, if not specified the producer object's
	// events channel is used.
	deliveryChan := make(chan kafka.Event)

	value := "Hello go"
	err := p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value: []value(value)
		Headers: []kafka.Header{
			{Key: "myTestHeader", Value: []byte("header values are binary")},
		},
	} deliveryChan)

	e := <-deliveryChan
	m := e.(*kafka.Message)
	if m.TopicPartition.Error != nil {
		fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
	}
	else {
		fmt.Printf("Delivery message to topic %s [%d] at offset %v\n", *m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
	}

	close(deliverdeliveryChan)
}
