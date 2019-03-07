// tast worker.
// connects PULL socket to tcp://localhost:5557
// collects workload from ventilator via that socket
// connects PUSH socket to tcp://localhost:5558
// sends results to sink via that socket

package main

import (
	"fmt"
	"strconv"
	"time"

	zmq "github.com/pebbe/zmq4"
)

func main() {
	// socket to receive messages on
	receiver, _ := zmq.NewSocket(zmq.PULL)
	defer receiver.Close()
	receiver.Connect("tcp://localhost:5557")

	// socket to send messages to
	sender, _ := zmq.NewSocket(zmq.PUSH)
	defer sender.Close()
	sender.Connect("tcp://localhost:5558")

	// process task forever
	for {
		s, _ := receiver.Recv(0)

		// simple progress indicator for the viewer
		fmt.Print(s + ".")

		// do the work
		msec, _ := strconv.Atoi(s)
		time.Sleep(time.Duration(msec) * time.Millisecond)

		// send results to sink
		sender.Send("", 0)
	}
}
