// go supports pointers.
// allowing you to pass references to values and records
// within your program.

package main

import "fmt"

// we'll show how pointers work in contrast to values with
// 2 functions: 'zeroval' and 'zeroptr'. 'zeroval' has an
// 'int' parameter, so arguments will be passed to it by
// balue. 'zeroval' will get a copy of 'ival' distinct
// from the one in the calling fuction.
func zeroval(ival int) {
	ival = 0
}

// 'zeroptr' in contrast has an '*int' parameter, meaning
// htat it takes an 'int' pointer. the '*ptr' code in the
// function body then dereferences the pointer from its
// memory address to the current value at that address.
// Assigning a value to a dereferenced pointer changes the
// value at the referenced address.
func zeroptr(iptr *int) {
	*iptr = 0
}

func main() {
	i := 1
	fmt.Println("initial: ", i)

	zeroval(i)
	fmt.Println("zeroval: ", i)

	// the '&i' syntax gives the memory address of 'i'.
	// i.e. a pointer to 'i'.
	zeroptr(&i)
	fmt.Println("zeroptr: ", i)

	// pointers can be printed too.
	fmt.Println("pointer: ", &i)
}
