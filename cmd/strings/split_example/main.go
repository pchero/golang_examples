package main

import (
	"fmt"
	"strings"
)

func main() {
	test1 := "tmp/facbb670-8d63-11ed-b89c-2f4bcc1326a1.zip"

	tmps := strings.Split(test1, "/")
	tmp := tmps[len(tmps)-1]

	fmt.Printf("res: %s\n", tmp)
}
