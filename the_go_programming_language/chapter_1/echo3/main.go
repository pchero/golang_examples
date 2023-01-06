package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Println(strings.Join(os.Args[1:], " "))
}

// $ go run cmd/the_go_programming_language/chapter_1/echo3/main.go test 1 2 3 4 5
// test 1 2 3 4 5
