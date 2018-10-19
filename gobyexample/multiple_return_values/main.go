// go has built-in support for _multiple return values_.
// this feature is used often in idiomatic go, for example
// to return both result and error values from a function

package main

import "fmt"

// the '(int int)' in this function signature shows that
// the function reutnrs 2 'int's.
func vals() (int, int) {
	return 3, 7
}

func main() {
	// here we use the 2 different return values from the
	// call with _multiple assignment_.
	a, b := vals()
	fmt.Println(a)
	fmt.Println(b)

	// if you only want a subset of the returned values,
	// use the blank identifier '_'.
	_, c := vals()
	fmt.Println(c)

	// this doesn't working.
	d := vals()
	fmt.Println(d)

}
