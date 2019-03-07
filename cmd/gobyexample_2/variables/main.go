// in go, variables are explicitly declared and used by
// the compiler to e.g. check type-correctness of function
// calls.

package main

import "fmt"

func main() {
	// 'var' declares 1 or more variables.
	var a = "initial"
	fmt.Println(a)

	// you can declare multiple variables at once.
	var b, c int = 1, 2
	fmt.Println(b, c)

	// to will infer the type of initialized variables.
	var d = true
	fmt.Println(d)

	// Variables declared without a corresponding
	// initialization are zero-valued. for example, the
	// zero value for an 'int' is '0'.
	var e int
	fmt.Println(e)

	// the ':=' syntax is shorthand for declaring and
	// initializing a variable, e.g. for
	// 'var f string = "short"' in this case.
	f := "short"
	fmt.Println(f)
}
