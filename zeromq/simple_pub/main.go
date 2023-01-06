//
//  Pathological publisher
//  Sends out 1,000 topics and then one random update per second
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
	publisher, err := zmq.NewSocket(zmq.PUB)
	if err != nil {
		fmt.Printf("Could not create a socket. err: %v", err)
	}

	if len(os.Args) == 2 {
		publisher.Bind(os.Args[1])
	} else {
		fmt.Printf("Binding the address.\n")
		if err := publisher.Bind("tcp://*:5555"); err != nil {
			// if err := publisher.Bind("inproc://simplepubsub"); err != nil {
			fmt.Printf("Could not bind a socket. err: %v", err)
		}
	}
	fmt.Printf("Creating socket\n")

	//  Ensure subscriber connection has time to complete
	time.Sleep(time.Second)

	//  Send out all 1,000 topic messages
	for topic_nbr := 0; topic_nbr < 1000; topic_nbr++ {
		_, err := publisher.SendMessage(fmt.Sprintf("%d", topic_nbr), "Save Roger")
		if err != nil {
			fmt.Println(err)
		}
	}
	fmt.Printf("Sent out topic messages\n")

	_, err = publisher.SendMessage(fmt.Sprintf("%d", 4321098), "Save Roger")

	//  Send one random update per second
	rand.Seed(time.Now().UnixNano())
	for {
		time.Sleep(time.Second)

		// topic := rand.Intn(1000)
		topic := 4321098
		// _, err := publisher.SendMessage(fmt.Sprintf("%03d", topic), "Off with his head!")
		_, err := publisher.SendMessage(fmt.Sprintf("%d", topic), "Off with his head!")

		if err != nil {
			fmt.Println(err)
		}

		fmt.Printf("Sent message. target: %d\n", topic)
	}
}
