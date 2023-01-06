package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/streadway/amqp"
)

var rabbitAddr = flag.String("rabbit_addr", "amqp://guest:guest@localhost:5672", "Rabbitmq address.")
var rabbitQueueName = flag.String("rabbit_queue_name", "rpc_test", "Queue name for RPC.")

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func rpcSendRecv() (res int, err error) {
	conn, err := amqp.Dial(*rabbitAddr)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
	fmt.Println("Connected to rabbitmq!")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()
	fmt.Println("Opend channel!")

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // noWait
		nil,   // arguments
	)
	fmt.Println("Queue declared!")
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")
	fmt.Println("Consumed!")

	corrID := randomString(32)

	err = ch.Publish(
		"",                        // exchange
		"asterisk_ari_10.164.0.3", // routing key
		false,                     // mandatory
		false,                     // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corrID,
			ReplyTo:       q.Name,
			Body:          []byte("request message"),
		})
	failOnError(err, "Failed to publish a message")
	fmt.Println("Published!")

	for d := range msgs {
		if corrID == d.CorrelationId {
			res, err = strconv.Atoi(string(d.Body))
			failOnError(err, "Failed to convert body to integer")
			break
		}
	}

	return
}

func main() {
	flag.Parse()
	rand.Seed(time.Now().UTC().UnixNano())

	// n := bodyFrom(os.Args)

	// log.Printf(" [x] Requesting fib(%d)", n)
	res, err := rpcSendRecv()
	failOnError(err, "Failed to handle RPC request")

	log.Printf(" [.] Got %d", res)
}

// func bodyFrom(args []string) int {
// 	var s string
// 	if (len(args) < 2) || os.Args[1] == "" {
// 		s = "30"
// 	} else {
// 		s = strings.Join(args[1:], " ")
// 	}
// 	n, err := strconv.Atoi(s)
// 	failOnError(err, "Failed to convert arg to integer")
// 	return n
// }
