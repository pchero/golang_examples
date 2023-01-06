package main

import (
	"flag"

	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

var rabbitAddr = flag.String("rabbit_addr", "amqp://guest:guest@localhost:5672", "rabbitmq service address.")
var rabbitQueue = flag.String("rabbit_queue", "example_delayed_queue", "rabbitmq delayed message example queue name.")

func main() {
	flag.Parse()
	logrus.SetLevel(logrus.DebugLevel)

	conn, err := amqp.Dial(*rabbitAddr)
	if err != nil {
		logrus.Errorf("Could not initiate the connection. err: %v", err)
		return
	}

	channel, err := conn.Channel()
	if err != nil {
		logrus.Errorf("Could not initiate the channel. err: %v", err)
		return
	}

	queue, err := channel.QueueDeclare("example_test", false, true, false, false, nil)
	if err != nil {
		logrus.Errorf("Could not initiate the queue. err: %v", err)
		return
	}

	// create a delayed exchange
	exchangeArgs := make(amqp.Table)
	exchangeArgs["x-delayed-type"] = "direct"
	if err := channel.ExchangeDeclare("example_delay", "x-delayed-message", false, true, false, false, exchangeArgs); err != nil {
		logrus.Errorf("Could not declare the exchange. err: %v", err)
		return
	}

	if err := channel.QueueBind(queue.Name, "bin-manager.call-manager", "example_delay", false, nil); err != nil {
		logrus.Errorf("Could not bind the queue. err: %v", err)
		return
	}

	// send a message
	headers := make(amqp.Table)
	headers["x-delay"] = 10000 // 10 sec
	channel.Publish("example_delay", "bin-manager.call-manager", false, false, amqp.Publishing{
		ContentType: "plain/text",
		Body:        []byte("hello world"),
		Headers:     headers,
	})

	messages, err := channel.Consume(
		"example_test",
		"tester",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logrus.Errorf("Could not consume the message. err: %v", err)
		return
	}

	for message := range messages {
		logrus.Infof("Received message: %s", string(message.Body))
	}

}
