// https://gobyexample.com/exit

// Use 'os.Exit' to immemdiately exit with a given
// status

package main

import (
	"fmt"
	"os"
)

func main() {
	// 'defer's will not be run when using 'os.Exit', so
	// this 'fmt.Println' will never be called.
	defer fmt.Println("!")

	// Exit with status 3
	os.Exit(3)
}

// Not that unlike e.g. C, Go does not use an integer
// return value from 'main' to indicate exit status. If
// you'd like to exit with a non-zero status you should
// use 'os.Exit'

// $ go run cmd/gobyexample_2/e
// xit/main.go
// exit status 3

// If you run exit.go using go run, the exit will be picked up
// by go and printed.

// By building and executing a binary you can see the status
// in the terminal.

// Note that the ! from this program never got printed.
