package main

import (
	"fmt"
)

func main() {
	// res := make([]string, 1)
	res := []string{}

	foundFiles := []string{
		"test1",
		"test2",
		"test3",
		"test4",
	}

	for _, file := range foundFiles {

		res = append(res, file)
	}

	fmt.Println(foundFiles)
	fmt.Println("-------------------------")
	fmt.Println(res)
}
