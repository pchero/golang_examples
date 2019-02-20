// go has built-in support for multiple return values.
// this feature is used often in idiomatic go, for example
// to return both result and error values from a function.

package main

import "fmt"

// the '(int, int)' int this function sigurature shows that
// the function returns 2 'int's.
func vals() (int, int) {
	return 3, 7
}

func main() {
	// here we use the 2 different return values from the
	// call with multiple assignment.
	a, b := vals()
	fmt.Println(a)
	fmt.Println(b)

	// if you only want a subset of the returned values,
	// use the blank identifier '_'.
	_, c := vals()
	fmt.Println(c)
}
