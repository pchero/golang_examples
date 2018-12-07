// tast worker.
// connects PULL socket to tcp://localhost:5557
// collects workload from ventilator via that socket
// connects PUSH socket to tcp://localhost:5558
// sends results to sink via that socket

package main

import (
	"fmt"

	zmq "github.com/pebbe/zmq4"
)

func main() {
	// socket to receive messages on
	receiver, _ := zmq.NewSocket(zmq.PULL)
	defer receiver.Close()
	receiver.Bind("tcp://*:5558")

	// process task forever
	for {
		s, _ := receiver.Recv(0)

		// simple progress indicator for the viewer
		fmt.Println(s)
	}
}
