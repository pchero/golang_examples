// a common requirement in program is getting the number
// of seconds, milliseconds, or nanoseconds since the
// unix epoch.
// here's how to do it in go.

package main

import "fmt"
import "time"

func main() {
	// use 'time.Now'  with 'Unix' or 'UnixNano' to get
	// elapsed time since the Unix epoch in seconds or
	// nanoseconds, respectively.
	now := time.Now()
	secs := now.Unix()
	nanos := now.UnixNano()
	fmt.Println(now)
	fmt.Println(secs)
	fmt.Println(nanos)

	// Note that there is no 'UnixMillis', so to get the
	// milliseconds since epoch you'll need to manually
	// divide from nanoseconds.
	millis := nanos / 1000000
	fmt.Println(secs)
	fmt.Println(millis)
	fmt.Println(nanos)

	// you can also convert integer seconds or nanoseconds
	// since the epoch into the corresponding 'time'.
	fmt.Println(time.Unix(secs, 0))
	fmt.Println(time.Unix(0, nanos))
}
