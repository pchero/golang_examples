package main

import (
	"fmt"
	"strings"
)

func main() {
	res := strings.Trim("+123456789", "+")
	fmt.Printf("res: %s\n", res)
}
