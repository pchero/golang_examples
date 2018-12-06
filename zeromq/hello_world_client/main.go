// hello world client.

package main

import (
	"fmt"

	zmq "github.com/pebbe/zmq4"
)

func main() {
	// socket to talk to server
	fmt.Println("Connecting to hello world server...")
	requester, _ := zmq.NewSocket(zmq.REQ)
	defer requester.Close()
	requester.Connect("tcp://localhost:5555")

	for request_nbt := 0; request_nbt < 10; request_nbt++ {
		// send hello
		msg := fmt.Sprintf("Hello %d", request_nbt)
		fmt.Println("Sending ", msg)
		requester.Send(msg, 0)

		// wait for reply
		reply, _ := requester.Recv(0)
		fmt.Println("Received ", reply)
	}
}
