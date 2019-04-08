// https://gobyexample.com/command-line-arguments

// Command-line arguments
// are common way to parameterize execution of programs.
// For example, 'go run hello.go' uses 'run' and
// 'hello.go' arguments to the 'go' program.

package main

import (
	"fmt"
	"os"
)

func main() {
	// 'os.Args' prvides access to raw command-line
	// arguments. Note that the first value in this slice
	// is the path to the program, and 'os.Args[1:]'
	// holds the arguments to the program.
	argsWithProg := os.Args
	argsWithoutProg := os.Args[1:]

	// You can get individual args with normal indexing.
	arg := os.Args[3]

	fmt.Println(argsWithProg)
	fmt.Println(argsWithoutProg)
	fmt.Println(arg)
}

// $ go run cmd/gobyexample_2/command-line_arguments/main.go test1 test2 test3 test4
// [/var/folders/9h/hzppr6wd0bs11w42rtfg3qqr0000gp/T/go-build394682576/b001/exe/main test1 test2 test3 test4]
// [test1 test2 test3 test4]
// test3
