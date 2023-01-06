package main

import (
	"fmt"
	"regexp"
)

func main() {
	regex := regexp.MustCompile("^\\+*([0-9]+)$")

	matched := regex.Match([]byte("+123"))
	fmt.Println(matched)

	subMatch := regex.FindSubmatch([]byte("+123"))
	fmt.Printf("%s\n", subMatch[1])

	subMatch2 := regex.FindSubmatch([]byte("12345"))
	fmt.Println(len(subMatch2))
	fmt.Printf("%s\n", subMatch2[0])

	subMatch = regex.FindSubmatch([]byte("test"))
	if subMatch != nil {
		fmt.Println(len(subMatch))
		fmt.Printf("%s\n", subMatch[0])
	}

	subMatch = regex.FindSubmatch([]byte("+123456789"))
	if subMatch != nil {
		fmt.Println(len(subMatch))
		fmt.Printf("%s\n", subMatch[1])
	}

	return
}
