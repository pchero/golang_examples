package main

import "fmt"

func main() {
	fmt.Printf("hello world\n")

	for i := 0; i < 3; i++ {
		defer fmt.Printf("end. %d\n", i)

		fmt.Printf("test: %d\n", i)
	}

	fmt.Printf("main end\n")
}
