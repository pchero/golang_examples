// echo2 prints its command-line arguments.
package main

import (
	"fmt"
	"os"
)

func main() {
	s, sep := "", ""
	for _, arg := range os.Args[1:] {
		s += sep + arg
		sep = " "
	}
	fmt.Println(s)
}

// $ go run cmd/the_go_programming_language/chapter_1/echo2/main.go test 1 2 3 4
// test 1 2 3 4
