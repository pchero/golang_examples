// go supports recursive functions
// here's a classic factorial example.

package main

import "fmt"

// this 'fact' function calls itself until it reaches the
// base case of 'fact(0)'
func fact(n int) int {
	if n == 0 {
		return 1
	}

	return n * (fact(n - 1))
}

func main() {
	fmt.Println(fact(7))
}
