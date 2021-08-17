package main

import (
	"fmt"
	"time"
)

var DefaultMaxPlaybackTime = 10 * time.Minute

func main() {
	fmt.Printf("Hello, playground. duration: %d", DefaultMaxPlaybackTime)

	go func() {
		for {
			fmt.Printf("Duration: %d\n", DefaultMaxPlaybackTime)
			time.Sleep(time.Second * 1)
		}
	}()
	time.Sleep(time.Second * 3)

	DefaultMaxPlaybackTime = 4 * time.Minute
	time.Sleep(time.Second * 3)
}
