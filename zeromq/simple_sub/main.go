//
//  Pathological subscriber
//  Subscribes to one random topic and prints received messages
//

package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	zmq "github.com/pebbe/zmq4"
)

func main() {
	subscriber, err := zmq.NewSocket(zmq.SUB)
	if err != nil {
		fmt.Printf("Could not create socket. err: %v", err)
	}

	if len(os.Args) == 2 {
		subscriber.Connect(os.Args[1])
	} else {
		fmt.Printf("Connecting to the socket. localhost:5556\n")
		if err := subscriber.Connect("tcp://localhost:5555"); err != nil {
			// if err := subscriber.Connect("inproc://simplepubsub"); err != nil {
			fmt.Printf("Could not connect to the socket. err: %v", err)
		}
	}
	fmt.Printf("Creating socket\n")

	rand.Seed(time.Now().UnixNano())
	// subscription := fmt.Sprintf("%03d", rand.Intn(1000))
	subscription := fmt.Sprintf("%d", 432)
	if err := subscriber.SetSubscribe(subscription); err != nil {
		fmt.Printf("Could not create subscribe. err: %v", err)
	}

	fmt.Printf("Subscription topic: %s\n", subscription)

	for {
		fmt.Printf("Receiving the message...\n")
		msg, err := subscriber.RecvMessage(0)
		if err != nil {
			fmt.Printf("Could not receive the message. err: %v", err)
			break
		}
		topic := msg[0]
		data := msg[1]
		if topic != subscription {
			fmt.Printf("Topic: %s", topic)
			// panic("topic != subscription")
		}
		fmt.Println(data)
	}
}
