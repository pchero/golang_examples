// timers are for when you want to do
// something once in the future. tickers are for when
// you want to do something repeatedly at regular
// intervals. here's an example of a ticker that ticks
// periodically until we stop it.

package main

import "time"
import "fmt"

func main() {
	// tickers are a similar mechanism to timers: a
	// channel that is sent values. here we'll use the
	// 'range' builtin on the channel to iterate over
	// the values as they arrive every 500 ms.
	ticker := time.NewTicker(500 * time.Millisecond)
	go func() {
		for t := range ticker.C {
			fmt.Println("Tick at", t)
		}
	}()

	// tickers can be stopped like timers. once a ticker
	// is stopped it won't receive any more values on its
	// channel. we'll stop ours after 1600ms
	time.Sleep(1600 * time.Millisecond)
	ticker.Stop()

	fmt.Println("ticker stopped.")
}
