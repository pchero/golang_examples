// go offers built-in support for regular expressions.
// here are some examples of common regex-related tasks.

package main

import (
	"bytes"
	"fmt"
	"regexp"
)

func main() {
	// this tests whether a pattern matches a string.
	match, _ := regexp.MatchString("p([a-z]+)ch", "peach")
	fmt.Println(match)

	match, _ = regexp.MatchString("p(.*)chddfs", "peach")
	fmt.Println("match result:", match)

	// aobve we used a string pattern directly, but for
	// other regexp tasks you'll need to 'Compile' an
	// optimized 'Regexp' struct.
	r, _ := regexp.Compile("p([a-z]+)ch")

	// many methods are available on these structs. here's
	// a match test like we saw earlier.
	fmt.Println(r.MatchString("peach"))

	// this finds the match for the regexp.
	fmt.Println(r.FindString("peach puch"))

	// this also finds the first match but returns the
	// start and end indexes for the match insted of the
	// matching text.
	fmt.Println(r.FindStringIndex("peach puch"))

	// the 'submatch' variants include information about
	// both the whole-pattern matches and the submatches
	// within those matches. for example this will return
	// information for both 'p([a-z]+)ch' and '([a-z]+)'.
	fmt.Println(r.FindStringSubmatch("peach punch"))

	// similarly this will return information about the
	// indexes of matches and submatches.
	fmt.Println(r.FindStringSubmatchIndex("peach punch"))

	// the 'All' variants are available for the other
	// functions we saw above as well.
	fmt.Println(r.FindAllString("peach punch pinch", -1))

	// these 'All' variants are available for the other
	// functions we saw above as well.
	fmt.Println(r.FindAllStringSubmatchIndex("peach punch pinch", -1))

	// providing a non-negative interger as the second
	// argument to these functions will limit the number
	// of matches.
	fmt.Println(r.FindAllString("peach punch pinch", 2))

	// our examples above had string arguments and used
	// names like 'MatchString'. We can also provide
	// '[]byte
	// arguments and drop 'String' from the
	// function name.
	fmt.Println(r.Match([]byte("peach")))

	// when creating constants with regular expressions
	// you can use the 'MustCompile' variation of
	// 'Compile'. A plain 'Compile' won't work for
	// constants because it has 2 return values.
	r = regexp.MustCompile("p([a-z]+)ch")
	fmt.Println(r)

	// the 'regexp' package can also be used to replace
	// subsets of strings with other values.
	fmt.Println(r.ReplaceAllString("a peach", "<fruit>"))

	// the 'Func' variant allows you to transform matched
	// text with a given function.
	in := []byte("a peach")
	out := r.ReplaceAllFunc(in, bytes.ToUpper)
	fmt.Println(string(out))
}
