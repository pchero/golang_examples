package main

import "fmt"

func main() {
	var s *string

	s = nil

	fmt.Printf("test: %s", *s)
}
