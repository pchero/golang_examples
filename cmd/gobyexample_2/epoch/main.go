// https://gobyexample.com/epoch

// A common requirement in program is getting the number
// of seconds, milliseconds, or nanoseconds since the
// unix epoch.
// Here's how to do it in Go.

package main

import (
	"fmt"
	"time"
)

func main() {
	// Use 'time.Now' with 'Unix' or 'UnixNano' to get
	// elapsed time since the Unix epoch in seconds or
	// nanoseconds, respectively.
	now := time.Now()
	secs := now.Unix()
	nanos := now.UnixNano()
	fmt.Println(now)

	// Note that there is no 'UnixMillis', so to get the
	// milliseconds since epoch you'll need to manually
	// divide from nanoseconds.
	millis := nanos / 1000000
	fmt.Println(secs)
	fmt.Println(millis)
	fmt.Println(nanos)

	// You can also convert integer seconds or nanoseconds
	// since the epoch into the corresponding 'time'.
	fmt.Println(time.Unix(secs, 0))
	fmt.Println(time.Unix(0, nanos))
}

// $ go run cmd/gobyexample_2/epoch/main.go
// 2019-04-02 08:55:13.860531 +0200 CEST m=+0.000376008
// 1554188113
// 1554188113860
// 1554188113860531000
// 2019-04-02 08:55:13 +0200 CEST
// 2019-04-02 08:55:13.860531 +0200 CEST
