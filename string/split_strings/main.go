package main

import (
	"fmt"
	"strings"
)

func main() {
	tmp := "test1|test2|test3"

	tmpArray := strings.Split(tmp, "|")

	fmt.Println(tmpArray)
}
