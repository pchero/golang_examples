// https://play.golang.org/p/I63ge2ISDs

// in a range example, we saw how 'for' and
// 'range' provide iteration over basic data structures.
// we can also use this syntax to iterate over
// values received from a channel.

// This example also showed that it's possible to close a non-
// empty channel but still have the remaining values be
// received.
package main

import (
	"fmt"
	"time"
)

func main() {
	// We'll iterate over 2 values in the 'queue' channel.
	queue := make(chan string, 2)
	queue <- "one"
	queue <- "two"

	// This 'range' iterates over each elements as it's
	// received from 'queue'. Because we 'close'd the
	// channel above, the iteration terminates after
	// receiving the 2 elements.
	go func() {
		for elem := range queue {
			fmt.Println(elem)
		}
	}()

	time.Sleep(1 * time.Second)

	queue <- "three"
	queue <- "four"

	close(queue)

	time.Sleep(1 * time.Second)

}
