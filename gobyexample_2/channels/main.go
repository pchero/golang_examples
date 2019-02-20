// channels are the pipes that connect concurrent
// goroutines. you can send values into channels from one
// goroutine and receive those vales into another
// goroutine.

package main

import "fmt"

func main() {
	// create a new channel with 'make(chan val-type)'.
	// channels are typed by the values they convey.
	messages := make(chan string)

	// send a value into a channel using the 'channel <-'
	// syntax. here we send "ping" to the 'messages'
	// channel we made above, from a new goroutine.
	go func() {
		messages <- "ping"
	}()

	// the '<-channel' syntax receives a value from the
	// channel. here we'll receive the "ping" message
	// we sent above and print it out.
	msg := <-messages
	fmt.Println(msg)
}
