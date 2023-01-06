// https://gobyexample.com/string-formatting

// Go offers  exellent support for string formatting in
// the 'printf' tradition. Here are some examples of
// common string formatting tasks.

package main

import (
	"fmt"
	"os"
)

type point struct {
	x int
	y int
}

func main() {
	// Go offers several printing "verbs" designed to
	// format general Go values. For example, this prints
	// an instance of our 'point' struct.
	p := point{1, 2}
	fmt.Printf("%v\n", p)

	// If the value is struct, the '%+v' variant will
	// include the struct's filed names
	fmt.Printf("%+v\n", p)

	// The '%#v' variant prints a Go syntax representation
	// of the value, i.e. the source code snippet that
	// would produce that value.
	fmt.Printf("%#v\n", p)

	// To print the type of a value, use '%T'.
	fmt.Printf("%T\n", p)

	// Formatting booleans is straight-forward.
	fmt.Printf("%t\n", true)

	// There are many options for formatting integers.
	// Use '%d' for standard, base-10 formatting.
	fmt.Printf("%d\n", 14)

	// This prints the character corresponding to the
	// given integer.
	fmt.Printf("%c\n", 33)

	// '%x' provides hex encoding.
	fmt.Printf("%x\n", 456)

	// There are also several formatting options for
	// floats. For basic decimal formatting use '%f'.
	fmt.Printf("%f\n", 78.9)

	// '%e' and '%E' format the float in (slightly
	// different version of) scientific notation.
	fmt.Printf("%e\n", 123400000.0)
	fmt.Printf("%E\n", 123400000.0)

	// For basic string printing use '%s'.
	fmt.Printf("%s\n", "\"string\"")

	// To double-quote string as in Go source, use '%q'.
	fmt.Printf("%q\n", "\"string\"")

	// As with integers seen earlier, '%x' renders
	// the string in base-16, with two output characters
	// per byte of input.
	fmt.Printf("%x\n", "hex this")

	// To print a representation of a pointer, use '%p'.
	fmt.Printf("%p\n", &p)

	// When formatting numbers you will often want to
	// control the width and precision of the resulting
	// figure. To specify the width of an integer, use a
	// number after the '%' in the verb. By default the
	// result will be right-justified and padded with
	// spaces.
	fmt.Printf("|%6d|%6d|\n", 12, 345)

	// You can also specify the width of printed floats,
	// though usually you'll also want to restrict the
	// decimal precision at the same time with the
	// width.precision syntax.
	fmt.Printf("|%6.2f|%6.2f|\n", 1.2, 3.45)

	// To left-justify, use the '-' flag.
	fmt.Printf("|%-6.2f|%-6.2f|\n", 1.2, 3.45)

	// You may also want to control width when formatting
	// strings, especially to ensure that they align in
	// table-like output. For basic right-justified width.
	fmt.Printf("|%6s|%6s|\n", "foo", "b")

	// You left-justify use the '-' flag as with numbers.
	fmt.Printf("|%-6s|%-6s|\n", "foo", "b")

	// So far we've seen 'Printf', which prints the
	// formatted string to 'os.Stdout'. 'Sprintf' formats
	// and returns a string without printing it anywhere.
	s := fmt.Sprintf("a %s", "string")
	fmt.Println(s)

	// You can format+print to 'io.Writers' other than
	// 'os.Stdout' using 'Fprintf'.
	fmt.Fprintf(os.Stderr, "an %s\n", "error")
}

// $ go run cmd/gobyexample_2/string_formatting/main.go
// {1 2}
// {x:1 y:2}
// main.point{x:1, y:2}
// main.point
// true
// 14
// !
// 1c8
// 78.900000
// 1.234000e+08
// 1.234000E+08
// "string"
// "\"string\""
// 6865782074686973
// 0xc000018090
// |    12|   345|
// |  1.20|  3.45|
// |1.20  |3.45  |
// |   foo|     b|
// |foo   |b     |
// a string
// an error
