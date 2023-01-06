package main

import "fmt"

func main() {
	test := []int{10, 9, 8, 7, 6}

	fmt.Printf("org value: %v", test)
	for i, v := range test {
		fmt.Printf("org value: %d\n", v)

		v = i
	}
	fmt.Printf("change value: %v", test)

}
