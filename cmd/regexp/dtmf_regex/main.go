package main

import (
	"fmt"
	"regexp"
)

func main() {

	strMatch := "[^0-9a-dA-D*#]"

	matched, err := regexp.Match(strMatch, []byte("123"))
	fmt.Println(matched, err)

	matched, err = regexp.Match(strMatch, []byte("abc"))
	fmt.Println(matched, err)

	matched, err = regexp.Match(strMatch, []byte("0123456789abcdABCD*#"))
	fmt.Println(matched, err)

	matched, err = regexp.Match(strMatch, []byte("123f"))
	fmt.Println(matched, err)

	matched, err = regexp.Match(strMatch, []byte(""))
	fmt.Println(matched, err)

}
