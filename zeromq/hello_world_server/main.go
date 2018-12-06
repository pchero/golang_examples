// hello world zeromq server

package main

import (
	"fmt"
	"time"

	zmq "github.com/pebbe/zmq4"
)

func main() {
	// socket to talk to client
	responder, _ := zmq.NewSocket(zmq.REP)
	defer responder.Close()

	responder.Bind("tcp://*:5555")

	for {
		// wait for next request from client
		msg, _ := responder.Recv(0)
		fmt.Println("Received ", msg)

		// do some work
		time.Sleep(time.Second)

		// send reply back to client
		reply := "World"
		responder.Send(reply, 0)
		fmt.Println("Sent ", reply)
	}

}
