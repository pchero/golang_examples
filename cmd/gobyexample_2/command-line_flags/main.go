// https://gobyexample.com/command-line-flags

// Command-line flags
// are a common way to specify options for command-line
// program. For example, in 'wc -l' the '-l' is a
// command-line flag.

package main

import (
	"flag"
	"fmt"
)

// Go provides a 'flag' package supporting basic
// command-line flag parsing. We'll use this package to
// implement our example command-line program.

func main() {
	// Basic flag declarations are available for string,
	// integer, and boolean options. Here we declare a
	// string flag 'word' with a default value "foo"
	// and a short description. This 'flag.String' function
	// returns a string pointer (not a string value);
	// we'll see how to use this pointer below.
	wordPtr := flag.String("word", "foo", "a string")

	// This declares 'numb' and 'fork' flags, using a
	// similar approach to the 'word' flag.
	numbPtr := flag.Int("numb", 42, "an int")
	boolPtr := flag.Bool("fork", false, "a bool")

	// It's also possible to declare an option that uses an
	// exisiting var declared elsewhere in the program.
	// Note that we need to pass in a pointer to the flag
	// declaration function.
	var svar string
	flag.StringVar(&svar, "svar", "bar", "a string var")

	// Once all flags are declared, call 'flag.Parse()'
	// to execute the command-line parsing.
	flag.Parse()

	// Here we'll just dump out the parsed options and
	// any trailing positional arguments. Note that we
	// need to dereference the pointers with e.g. '*wordPtr'
	// to get the actual option values.
	fmt.Println("word: ", *wordPtr)
	fmt.Println("numb: ", *numbPtr)
	fmt.Println("fork: ", *boolPtr)
	fmt.Println("svar: ", svar)
	fmt.Println("tail: ", flag.Args())
}

// $ go run cmd/gobyexample_2/command-line_flags/main.go
// word:  foo
// numb:  42
// fork:  false
// svar:  bar
// tail:  []

// $ go run cmd/gobyexample_2/command-line_flags/main.go -h
// Usage of /var/folders/9h/hzppr6wd0bs11w42rtfg3qqr0000gp/T/go-build803906378/b001/exe/main:
//   -fork
//         a bool
//   -numb int
//         an int (default 42)
//   -svar string
//         a string var (default "bar")
//   -word string
//         a string (default "foo")
// exit status 2

// $ go run cmd/gobyexample_2/command-line_flags/main.go -fork true -word hello world
// word:  foo
// numb:  42
// fork:  true
// svar:  bar
// tail:  [true -word hello world]
