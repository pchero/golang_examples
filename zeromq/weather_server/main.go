// weather update server
// bind pub socket to tcp://*:5556
// publishes random weather updates

package main

import (
	"fmt"
	"math/rand"
	"time"

	zmq "github.com/pebbe/zmq4"
)

func main() {
	// prepare our publisher
	publisher, _ := zmq.NewSocket(zmq.PUB)
	defer publisher.Close()

	publisher.Bind("tcp://*:5556")
	publisher.Bind("ipc://weather.ipc")

	// initialize random number generator
	rand.Seed(time.Now().UnixNano())

	// loop for a while apparently
	for {
		// get values that will fool the boss
		zipcode := rand.Intn(100000)
		temperature := rand.Intn(215) - 80
		relhumidity := rand.Intn(50) + 10

		// send message to all subscribers
		msg := fmt.Sprintf("%05d %d %d", zipcode, temperature, relhumidity)
		publisher.Send(msg, 0)
	}
}
