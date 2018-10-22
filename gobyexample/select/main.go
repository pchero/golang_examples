// go's select lets you wait on multiple channel
// operations. combining goroutines and channels with
// select is a powerful feature of go

package main

import "time"
import "fmt"

func main() {
	// for our example we'll select across two channels.
	c1 := make(chan string)
	c2 := make(chan string)

	// each channel will receive a value after some amount
	// of time, to simulate e.g. blocking RPC oerpations
	// executing in concurrent goroutines.
	go func() {
		time.Sleep(1 * time.Second)
		c1 <- "one"
	}()

	go func() {
		time.Sleep(2 * time.Second)
		c2 <- "two"
	}()

	// we'll use 'select' to await both of these values
	// simultaneously, printing each one as it arrives.
	for i := 0; i < 2; i++ {
		select {
		case msg1 := <-c1:
			{
				fmt.Println("received", msg1)
			}
		case msg2 := <-c2:
			{
				fmt.Println("received", msg2)
			}
		}
	}
}
