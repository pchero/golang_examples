// 'for' is go's only looping construct. here are
// three basic types of 'for' loops.

package main

import "fmt"

func main() {
	// the most basic type, with a single condition.
	i := 1
	for i <= 3 {
		fmt.Println(i)
		i = i + 1
	}

	// a classic initial/condition/after 'for' loop.
	for j := 7; j <= 9; j++ {
		fmt.Println(j)
	}

	// 'for' without a condition will loop repeatedly
	// until you 'break' out of the loop or 'return' from
	// the eclosing function.
	for {
		fmt.Println("loop")
		break
	}

	// you can also 'continue' to the next iteration of
	// the loop.
	for n := 0; n <= 5; n++ {
		if n%2 == 0 {
			continue
		}
		fmt.Println(n)
	}
}
