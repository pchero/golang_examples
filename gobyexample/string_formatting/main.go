// go offers excellent support for string formatting in
//the 'printf' tradition. here are some examples of
// common string formatting tasks.

package main

import (
	"fmt"
	"os"
)

type point struct {
	x, y int
}

func main() {
	// go offers several printing "verbs" designed to
	// format general go values. for example, this prints
	// an instance of our 'point' struct.
	p := point{1, 2}
	fmt.Printf("%v\n", p)

	// if the value is a struct, the '%+v' variant will
	// include the struct's filed names.
	fmt.Printf("%+v\n", p)

	// the '%#v' variant prints a go syntax representation
	// of the value, i.e. the source code snippet that
	// would produce that value.
	fmt.Printf("%#v\n", p)

	// to print the type of a value, use '%T'.
	fmt.Printf("%T\n", p)

	// formatting booleans is strait-forward.
	fmt.Printf("%t\n", true)

	// there are many options for formatting integers.
	// use '%d' for standard, base-10 formatting.
	fmt.Printf("%d\n", 123)

	// this prints a binary representation.
	fmt.Printf("%b\n", 14)

	// this prints the character corresponding to the
	// given integer
	fmt.Printf("%c\n", 33)

	// '%x' provides hex encoding.
	fmt.Printf("%x\n", 456)

	// there are also several formatting options for
	// floats. for basic decimal formatting use '%f'.
	fmt.Printf("%f\n", 78.9)

	// '%e' and '%E' format the float in (slightly
	// deifferent version of) scientific notation.
	fmt.Printf("%e\n", 12340000.0)
	fmt.Printf("%E\n", 12340000.0)

	// for basic string printing use '%s'.
	fmt.Printf("%s\n", "\"string\"")

	// to double-quote trings as in go source, use '%q'.
	fmt.Printf("%q\n", "\"string\"")

	// as with integers seen earlier, '%x' renders
	// the string in base-16, with two output characters
	// per byte of input.
	fmt.Printf("%x\n", "hex this")

	// to print a representation of a pointer, use '%p'.
	fmt.Printf("%p\n", &p)

	// when formatting number you will often want to
	// control the width and precision of the resulting
	// figure. to specify the width of an integer, use a
	// number after the '%' in the verb. by default the
	// result will be rith-justified and padded with
	// spaces.
	fmt.Printf("|%6d|%6d|\n", 12, 345)

	// you can also specify the width of printed floats,
	// though usually you'll also want to restrict the
	// decimal precision at the same time with the
	// width.precision syntax.
	fmt.Printf("|%6.2f|%6.2f|%6.2f|\n", 1.2, 3.45, 6.781)

	// to left-justify, use the '-' flag.
	fmt.Printf("|%-6.2f|%-6.2f|%-6.2f|\n", 1.2, 3.45, 6.789)

	// you may also want to control width when formatting
	// strings, especially to ensure that they align in
	// take-like output. for basic right-justified width.
	fmt.Printf("|%6s|%6s|%6s|\n", "foo", "b", "longer than 6")

	// to left-justify use the '-' flag as with numbers.
	fmt.Printf("|%-6s|%-6s|%-6s|\n", "foo", "b", "longer than 6")

	// so far we've seen 'Printf', which prints the
	// formatted string to 'os.Stdout'. 'Sprintf formats
	// and returns a string without printing it anywhere.
	s := fmt.Sprintf("a %s", "string")
	fmt.Println(s)

	// you can format+print to 'io.Writers' other than
	// 'os.Stdout' using 'Fprintf'.
	fmt.Fprintf(os.Stderr, "an %s\n", "error")
}
