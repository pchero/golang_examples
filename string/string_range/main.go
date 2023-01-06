package main

import "fmt"

func main() {
	tmp := "Hello world"

	for _, val := range tmp {
		fmt.Printf("%c", val)
	}
}
