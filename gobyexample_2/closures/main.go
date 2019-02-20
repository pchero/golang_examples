// go supports anonymous functions
// which can form closure.
// anonymous functions are useful when you want to define
// a function inline without having to name it.

package main

import "fmt"

// this function 'intSeq' returns another function, which
// we define anonymous in the body of 'intSeq'. the
// returned function closes over the variable 'i' to
// form a closure.
func intSeq() func() int {
	i := 0
	return func() int {
		i++
		return i
	}
}

func main() {
	// we call 'intSeq', assigning the result (a function)
	// to 'nextInt'. This function value captures its
	// own 'i' value, which will be updated each time
	// we call 'nextInt'
	nextInt := intSeq()

	// see the effect of the closure by calling 'nextInt'
	// a few times.
	fmt.Println(nextInt())
	fmt.Println(nextInt())
	fmt.Println(nextInt())

	// to confirm that the state is unique to that
	// paricular function, create and test a new one.
	newInts := intSeq()
	fmt.Println(newInts())

	fmt.Println(nextInt())

}
