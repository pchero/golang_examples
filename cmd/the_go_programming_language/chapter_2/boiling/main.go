// Boiling prints the boiling point of water.
package main

import "fmt"

const boilingF = 212.0

func main() {
	var f = boilingF
	var c = (f - 32) * 5 / 9

	fmt.Printf("boiling point = %gF or %gC\n", f, c)
}

// $ go run cmd/the_go_programming_language/chapter_2/boiling/main.go
// boiling point = 212F or 100C
