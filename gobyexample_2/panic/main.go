// https://gobyexample.com/panic

// A 'panic' typically means something went unexpectedly
// wrong. Mostly we use it to fail fast on errors that
// shouldn't occur during normal operation, or that we
// aren't prepared to handle gracefully.

package main

import (
	"fmt"
	"os"
)

func main() {

	// We'll use panic throughout this site to check for
	// unexpected errors. This is the only program on the
	// site designed to panic.
	panic("a problem")

	// A common use of panic is to abort if a function
	// returns an error value that we don't know how to
	// (or want to) handle. Here's an example of
	// 'panic'king if we get an unexpected error when creating
	// a new file.
	_, err := os.Create("/tmp/file")
	if err != nil {
		panic(err)
	}

	fmt.Println("test gogo")
}

// $ go run panic/main.go
// panic: a problem

// goroutine 1 [running]:
// main.main()
//         /Users/pchero/go/src/github.com/pchero/golang_examples/cmd/gobyexample_2/panic/main.go:20 +0x39
// exit status 2

// Note that unlike some languages which use exceptions for
// handling of many errors, in Go it's idiomatic to use error-
// indicating return values wherever possible.
