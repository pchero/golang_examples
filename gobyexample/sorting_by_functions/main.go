// sometimes we'll want to sort a collection by something
// other than its natural order. for example, suppose we
// wanted to sort strings by their length instead of
// alphabetically. here's an example of custom sorts
// in go

package main

import (
	"fmt"
	"sort"
)

// in order to sort by a custom function in go, we need a
// corresponding type. here we've created a 'byLength'
// type that is just an alias for the builtin '[]string'
// type
type byLength []string

// we implement 'sort.Interface' - 'Len', 'Less', and
// 'Swap' - on our type so we can use the 'sort' package's
// generic 'Sort' function. 'Len' and 'Swap'
// will usually be similar across type and 'Less' will
// hold the actual custom sorting logic. in our case we
// want to sort in order of increasing string length, so
// we use 'len(s[i])' and 'len(s[j])' here.
func (s byLength) Len() int {
	return len(s)
}
func (s byLength) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byLength) Less(i, j int) bool {
	return len(s[i]) < len(s[j])
}

// with all of this in place, we can now implement our
// custom sort casting the original 'fruits' slice to
// 'byLength', and then use 'sort.Sort' on that typed
// slice.
func main() {
	fruits := []string{"peach", "banana", "kiwi"}
	sort.Sort(byLength(fruits))
	fmt.Println(fruits)
}
