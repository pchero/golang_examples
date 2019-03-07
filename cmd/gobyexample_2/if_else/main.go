// branching with 'if' and 'else' in go is
// straight-forward.

package main

import "fmt"

func main() {
	// here's basic example.
	if 7%2 == 0 {
		fmt.Println("7 is even")
	} else {
		fmt.Println("7 is odd")
	}

	// you can have an 'if' statement without an else.
	if 8%4 == 0 {
		fmt.Println("8 is divisible by 4")
	}

	// a statement can precede conditionals; any variables
	// declared in this statement are available in all
	// branches.
	if num := 9; num < 0 {
		fmt.Println(num, "is negative")
	} else if num < 10 {
		fmt.Println(num, "has 1 digit")
	} else {
		fmt.Println(num, "has multiple digits")
	}
}

// note that you don't need parentheses around conditions
// in go, but that the braces are required
