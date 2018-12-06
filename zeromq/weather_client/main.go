package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	zmq "github.com/pebbe/zmq4"
)

func main() {
	// socket to tlak to server
	fmt.Println("Collecting updates from weather server...")
	subscriber, _ := zmq.NewSocket(zmq.SUB)
	defer subscriber.Close()
	subscriber.Connect("tcp://localhost:5556")

	// subscribe to zipcode, default is NYC, 10001
	filter := "10001 "
	if len(os.Args) > 1 {
		filter = os.Args[1] + " "
	}
	subscriber.SetSubscribe(filter)

	// process 100 updates
	total_temp := 0
	update_nbr := 0
	for update_nbr < 100 {
		msg, _ := subscriber.Recv(0)

		if msgs := strings.Fields(msg); len(msgs) > 1 {
			if temperature, err := strconv.Atoi(msgs[1]); err == nil {
				total_temp += temperature
				update_nbr++
			}
		}
	}

	fmt.Printf("Average temperature for zipcode '%s' was %dF\n\n", strings.TrimSpace(filter), total_temp/update_nbr)
}
