package main

import (
	"fmt"
	"time"
)

func main() {
	curTime := time.Now()

	fmt.Println(curTime)
	fmt.Println(curTime.Unix())
	fmt.Println(curTime.UnixNano())
}
