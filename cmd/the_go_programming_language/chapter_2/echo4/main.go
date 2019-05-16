// echo4 prints its command-line arguments.
package main

import (
	"flag"
	"fmt"
	"strings"
)

var n = flag.Bool("n", false, "omit trailing newline")
var sep = flag.String("s", " ", "separator")

func main() {
	flag.Parse()
	fmt.Print(strings.Join(flag.Args(), *sep))
	if !*n {
		fmt.Println()
	}
}

// go run cmd/the_go_programming_language/chapter_2/echo4/main.go -h
// Usage of /var/folders/9h/hzppr6wd0bs11w42rtfg3qqr0000gp/T/go-build946069622/b001/exe/main:
//   -n    omit trailing newline
//   -s string
//         separator (default " ")
// exit status 2
