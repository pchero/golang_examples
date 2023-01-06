// https://gobyexample.com/string-functions

// The standard library's 'string' package provides many
// useful string-related functions. Here are some examples
// to give you a sense of the package.

package main

import (
	"fmt"
	"strings"
)

// We alias 'fmt.Println' to a shorter name as we'll use
// it a lot below.
var p = fmt.Println

func main() {
	// Here's a sample of the functions available in
	// 'strings'. Since these are functions from the
	// package, not methods on the string object itself,
	// we need pass the string in question as the first
	// argument to the function. You can find more
	// functions in the 'strings'
	// package docs.
	p("Contains:  ", strings.Contains("test", "es"))
	p("Count:     ", strings.Count("test", "t"))
	p("HasPrefix: ", strings.HasPrefix("test", "te"))
	p("HasSuffix: ", strings.HasSuffix("test", "st"))
	p("Index:     ", strings.Index("test", "e"))
	p("Join:      ", strings.Join([]string{"a", "b", "c"}, "-"))
	p("Repeat:    ", strings.Repeat("a", 5))
	p("Replace    ", strings.Replace("foo", "o", "0", -1))
	p("Replace    ", strings.Replace("foo", "o", "0", 1))
	p("Replace    ", strings.Replace("foo", "o", "0", 0))
	p("Split:     ", strings.Split("a-b-c-d-e-f", "-"))
	p("ToLower:   ", strings.ToLower("TEST"))
	p("ToUpper:   ", strings.ToUpper("test"))
	p()

	// Not part of 'strings', but worth mentioning here, are
	// the mechanisms for getting the length of a string in
	// bytes and getting a byte by index.
	p("Len:       ", len("hello"))
	p("Char:      ", "hello"[1])
}

// Note that 'len' and indexing above work at the byte level.
// Go uses UTF-8 encoded strings, so this is often useful
// as-is. If you're working with potentially multi-byte
// characters you'll want to use encoding-aware operations.
// See https://blog.golang.org/strings for more information

// $ go run cmd/gobyexample_2/string_functions/main.go
// Contains:   true
// Count:      2
// HasPrefix:  true
// HasSuffix:  true
// Index:      1
// Join:       a-b-c
// Repeat:     aaaaa
// Replace     f00
// Replace     f0o
// Replace     foo
// Split:      [a b c d e f]
// ToLower:    test
// ToUpper:    TEST

// Len:        5
// Char:       101
