// primary mechanism for managing state in go is
// communication over channels. we saw this for example
// with worker pools. there are a few other
// options for managing state though. here we'll
// look at using the 'sync/atomic' package for atomic
// counters accessed by multiple goroutines.

package main

import (
	"fmt"
	"sync/atomic"
	"time"
)

func main() {
	// we'll use an unsigned integer to represent our
	// (always-positive) counter.
	var ops uint64

	// to simulate concurrent update, we'll start 50
	// goroutines that each increment the counter about
	// once a millisecond.
	for i := 0; i < 50; i++ {
		go func() {
			for {
				// to automatically increment the counter we
				// use 'AddUint64', giving it the memory
				// address of our 'ops' counter with the
				// '&' syntax.
				atomic.AddUint64(&ops, 1)

				// wait a bit between increments.
				time.Sleep(time.Millisecond)
			}
		}()
	}

	// wait a second to allow some ops to accumulate.
	time.Sleep(time.Second)

	// in order to safely use the counter while it's still
	// being updated by other goroutines, we extract a
	// copy of the current value into 'opsFinal' via
	// 'LoadUint64'. As above we need to give this
	// function the memory address '&ops' from which to
	// fetch the value.
	opsFinal := atomic.LoadUint64(&ops)
	fmt.Println("ops:", opsFinal)
}
