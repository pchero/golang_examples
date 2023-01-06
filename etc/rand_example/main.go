package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	min := 10000
	max := 11000
	fmt.Println(rand.Intn(max-min+1) + min)
}
