// go supports constants of character, string, boolean,
// and numeric values.

package main

import (
	"fmt"
	"math"
)

// 'const' declares a constant value.
const s string = "constant"

func main() {
	fmt.Println(s)

	// a 'const' statement can appear anywhere a 'var'
	// statement can.
	const n = 50000000

	// const expressions perform arithmetic with
	// arbitrary precision.
	const d = 3e20 / n
	fmt.Println(d)

	// a numeric constant has no type until it's given
	// one, such as by an explicit cast.
	fmt.Println(int64(d))

	// a number can be given a type by using it in a
	// context that requires one, such as a variable
	// assignment or function call. for example, here
	// 'math.Sin' expects a 'float64'.
	fmt.Println(math.Sin(n))
}
