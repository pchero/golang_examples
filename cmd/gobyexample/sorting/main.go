// go's 'sort' package implements sorting for builtins
// and user-defined types. we'll look at sorting for
// builtins first.

package main

import "fmt"
import "sort"

func main() {
	// sort methods are specific to the builtin type;
	// here's an example for strings. note that sorting is
	// in-place, so it changes the given slice and doesn't
	// return a new one.
	strs := []string{"c", "a", "b"}
	sort.Strings(strs)
	fmt.Println("Strings:", strs)

	// an example of sorting 'int's.
	ints := []int{7, 2, 4}
	sort.Ints(ints)
	fmt.Println("Ints:   ", ints)

	// we can also use 'sort' to check if a slice is
	// already in sorted order.
	s := sort.IntsAreSorted(ints)
	fmt.Println("Sorted: ", s)
}
