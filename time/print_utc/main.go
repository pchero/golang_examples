package main

import (
	"fmt"
	"strings"
	"time"
)

func main() {

	now := time.Now()

	fmt.Println(now)
	fmt.Println(now.UTC())

	nowTime := time.Now().UTC().String()
	res := strings.TrimSuffix(nowTime, " +0000 UTC")
	fmt.Println(res)
}
