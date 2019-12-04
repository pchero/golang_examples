package main

import (
	"fmt"
	"strconv"
)

func main() {
	n, err := strconv.ParseInt("32", 10, 64)
	if err != nil {
		fmt.Errorf("error: %w\n", err)
	}
	fmt.Printf("result: %d\n", n)

	t1, err := strconv.Atoi("test")
	if err != nil {
		fmt.Errorf("error: %w\n", err)
	}
	fmt.Printf("result: %d\n", t1)
}
