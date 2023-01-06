// when using channel as function parameters, you can
// specify if a channel is meant to only send or receive
// values. this specificity increases the type-safety or
// the program.

package main

import (
	"fmt"
)

// this "ping" function only accepts a channel for sending
// values. it would be a compile-time error to try to
// receive on this channel.
func ping(pings chan<- string, msg string) {
	pings <- msg
}

// the "pong" function accepts one channel for receives
// ("pings") and a second for sends ("pongs").
func pong(pings <-chan string, pongs chan<- string) {
	msg := <-pings
	pongs <- msg
}

func main() {
	pings := make(chan string, 1)
	pongs := make(chan string, 1)
	ping(pings, "passed message")
	pong(pings, pongs)
	fmt.Println(<-pongs)
}
