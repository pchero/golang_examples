// task ventilator
// binds PUSH socket to tcp://localhost:5557
// sends batch of tasks to worker via that socket

package main

import (
	"fmt"
	"math/rand"
	"time"

	zmq "github.com/pebbe/zmq4"
)

func main() {
	// socket to send messages on
	sender, _ := zmq.NewSocket(zmq.PUSH)
	defer sender.Close()

	// bind
	if err := sender.Bind("tcp://*:5557"); err != nil {
		fmt.Println("Could not bind to socket. err: ", err)
		return
	}

	// socket to send start of batch message on
	sink, _ := zmq.NewSocket(zmq.PUSH)
	defer sink.Close()
	if err := sink.Connect("tcp://localhost:5558"); err != nil {
		fmt.Println("Could not connect to socket. err: ", err)
		return
	}

	//
	fmt.Println("Press Enter when the workers are ready: ")
	var line string
	fmt.Scanln(&line)
	fmt.Println("Sending tasks to worker...")

	// the first message is "0" and signals start of batch
	sink.Send("0", 0)

	// initialize random number generator
	rand.Seed(time.Now().UnixNano())

	// send 100 taks
	total_msec := 0
	for task_nbr := 0; task_nbr < 100; task_nbr++ {
		// random workload from 1 to 100msecs
		workload := rand.Intn(100) + 1
		total_msec += workload
		s := fmt.Sprintf("%d", workload)
		sender.Send(s, 0)
	}
	fmt.Println("Total expected cost:", time.Duration(total_msec)*time.Millisecond)
	time.Sleep(time.Second) // give 0MQ time to deliver

}
