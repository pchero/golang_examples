// in a range example, we saw how 'for' and
// 'range' provide iteration over basic data structures.
// we can also use this syntax to iterate over
// values received from a channel.

package main

import "fmt"

func main() {
	// we'll iterate over 2 values in the 'queue' channel.
	queue := make(chan string, 2)
	queue <- "one"
	queue <- "two"
	close(queue)

	// this 'range' iterate over  element as it's
	// received from 'queue'. because we 'close'd the
	// channel aobve, the iteration terminates after
	// receiving the 2 elements.
	for elem := range queue {
		fmt.Println(elem)
	}
}
