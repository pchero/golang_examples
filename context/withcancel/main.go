// this example demonstrates the use of a cancelable context to prevent a gorouine leak.
// by the end of the example function, the goroutine started by gen will return without leaking.
package main

import (
	"fmt"

	"golang.org/x/net/context"
)

func main() {
	// gen generate integers in a separate gorouine and
	// sends them to the returned channel.
	// The callers of gen nedd to cancel the context once
	// they are done dunsuming generated integers not to leak
	// the internal goroutine started by gen.
	gen := func(ctx context.Context) <-chan int {
		dst := make(chan int)
		n := 1
		go func() {
			for {
				select {
				case <-ctx.Done():
					return // returning not to leak the goroutine.

				case dst <- n:
					n++
				}
			}
		}()

		return dst
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // cancel with we are finished consuming integers

	for n := range gen(ctx) {
		fmt.Println(n)
		if n == 5 {
			break
		}
	}
}
