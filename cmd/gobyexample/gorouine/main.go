// a goroutine is a lightweight thread of execution.

package main

import "fmt"

func f(from string) {
	for i := 0; i < 3; i++ {
		fmt.Println(from, ":", i)
	}
}

func main() {
	// suppose we have a function call 'f(s)'. here's how
	// we'd call that in the usual way, running it
	// synchronously.
	f("direct")

	// to invoke this function in a goroutine, use
	// 'go f(s)'. this new goroutine will execute
	// concurrently with the calling one.
	go f("goroutine")

	// you can also start a goroutine for an anonymous
	// function call.
	go func(msg string) {
		fmt.Println(msg)
	}("going")

	// our two function calls are running asynchronously in
	// separate goroutines now, so execution falls through
	// to here. this 'scanln' requires we press a key before
	// the program exists.
	fmt.Scanln()
	fmt.Println("done")
}
