// a 'panic' typically means something went unexpectedly
// wrong. mostly we use it to fail fast on errors that
// shouldn't occur during normal operation, or that we
// aren't prepared to handle gracefully.

package main

import "os"

func main() {
	// we'll use panic throughout this site to check for
	// unexpected errors. this is the only program on the
	// site designed to panic.
	panic("a problem")

	// a common use of panic is to abort if a function
	// returns an error value that we don't know how to
	// (or want to) handle. here's an example of
	// 'panic'king if we get an unexpected error when creating a new file.
	_, err := os.Create("/tmp/file")
	if err != nil {
		panic(err)
	}
}
