// https://gobyexample.com/sorting

// Go's 'sort' package implements sorting for builtins
// and user-defined types. We'll look at sorting for
// builtins first.

package main

import (
	"fmt"
	"sort"
)

func main() {
	// Sort methods are specific to the builtin type;
	// here's an example for strings. Note that sorting is
	// in-place, so it changes the give slice and doesn't
	// return new one.
	strs := []string{"c", "a", "b"}
	fmt.Println("Before sorting string:", strs)

	sort.Strings(strs)
	fmt.Println("After sorting string:", strs)

	// An example of sorting 'int's.
	ints := []int{7, 2, 4}
	fmt.Println("Before sorting int: ", ints)

	sort.Ints(ints)
	fmt.Println("After soring int:", ints)

	// We can also use 'sort' to check if a slice is
	// already in sorted order.
	s := sort.IntsAreSorted(ints)
	fmt.Println("Sorted:", s)

}

// $ go run cmd/gobyexample_2/sorting/main.go
// Before sorting string: [c a b]
// After sorting string: [a b c]
// Before sorting int:  [7 2 4]
// After soring int: [2 4 7]
// Sorted: true

// Running this program prints string and int
// slices and true as the result of AreSorted test.
