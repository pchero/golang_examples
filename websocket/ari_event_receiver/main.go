package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pebbe/zmq4"
)

var addr = flag.String("addr", "localhost:8088", "http service address")
var account = flag.String("account", "asterisk:asterisk", "Asterisk ari user info. id:password")
var subscribeAll = flag.String("subscribe_all", "true", "Subscribe all")
var application = flag.String("application", "event-reciver", "Application name")
var destination = flag.String("destination", "tcp://localhost:5558", "Destination zmq socket address")

// create channel
var messages = make(chan string, 1024000)

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// asterisk message receiver
	go recvEventFromAst()

	// message publisher
	go publishEvent()

	// interrupt handler
	for {
		select {
		case <-interrupt:
			log.Println("Interrupted.")

			// todo:
			// how to close all go routines?

			return
		}
	}

	fmt.Println("Finished.")

	return
}

func recvEventFromAst() {
	// create url query
	rawquery := fmt.Sprintf("api_key=%s&subscribeAll=%s&app=%s", *account, *subscribeAll, *application)

	u := url.URL{
		Scheme:   "ws",
		Host:     *addr,
		Path:     "/ari/events",
		RawQuery: rawquery,
	}
	log.Printf("Dial string: %s", u.String())

	for {
		// connect
		c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			log.Println("Could not connect to server. err: ", err)

			// sleep for every second
			time.Sleep(1 * time.Second)
			continue
		}
		defer c.Close()

		// receiver
		for {
			msgType, msgStr, err := c.ReadMessage()
			if err != nil {
				log.Printf("Could not read message. msgType: %d, err: %s", msgType, err)
				break
			}
			log.Printf("Message received. type: %d, message: %s", msgType, msgStr)

			// insert msg into queue
			messages <- string(msgStr)
		}

		// sleep 1 second for reconnect
		time.Sleep(1 * time.Second)
	}
}

func publishEvent() {
	sender, _ := zmq4.NewSocket(zmq4.PUSH)
	defer sender.Close()

	err := sender.Connect(*destination)
	if err != nil {
		fmt.Println("Could not connect to destination. err: ", err)
		return
	}

	for {

		select {
		case msg := <-messages:
			fmt.Println("Sending message: ", msg)
			sender.Send(msg, 0)
		}
	}
}
