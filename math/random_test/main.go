package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	// s3 := rand.NewSource(42)
	s3 := rand.NewSource(time.Now().UnixNano())

	r3 := rand.New(s3)
	fmt.Print(r3.Intn(100), ",")
	fmt.Print(r3.Intn(100))
	fmt.Println()

	fmt.Println("%d", rand.Reader)
}
