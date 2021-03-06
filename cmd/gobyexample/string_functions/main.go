// the standard library's 'string' package provides many
// useful string-related functions. here are some examples
// to give you a sense of the package.

package main

import s "strings"
import "fmt"

// we alais 'fmt.Println' to a shorter name as we'll use
// it a lot below.
var p = fmt.Println

func main() {
	// here's a sample of the functions available in
	// 'strings'. since these are functions from the
	// pacakge, not methods on the string obejct itself,
	// we need pass the string in question as the first
	// argument to the function. you can find more
	// functions in the 'string' package docs.
	p("Contains:    ", s.Contains("test", "es"))
	p("Count:       ", s.Count("test", "t"))
	p("HashPrefix:  ", s.HasPrefix("test", "te"))
	p("HashSuffix:  ", s.HasSuffix("test", "st"))
	p("Index:       ", s.Index("test", "e"))
	p("Join:        ", s.Join([]string{"a", "b"}, "-"))
	p("Repeat:      ", s.Repeat("a", 5))
	p("Replace:     ", s.Replace("foo", "o", "0", -1))
	p("Replace:     ", s.Replace("foo", "o", "0", 1))
	p("Split:       ", s.Split("a-b-c-d-e", "-"))
	p("ToLower:     ", s.ToLower("TEST"))
	p("ToUpper:     ", s.ToUpper("test"))
	p()

	// not part of 'strings', but worth metioning here, are
	// the mechanisms for getting the length of a string in
	// bytes and getting a byte by index.
	p("Len:         ", len("hello"))
	p("Char:        ", "hello"[1])
}

// note that 'len' and indexing above work at the byte level.
// go uses UTF-8 encoded strings, so this is often useful
// as-is. if you're working with potentially multi-byte
// characters you'll want to use encoding-aware operations.
