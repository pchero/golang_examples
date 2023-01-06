package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	fmt.Printf("hello world\n")

	// regex := regexp.MustCompile("^\\+*0*([1-9][0-9]*)$")
	regex := regexp.MustCompile(`^[\+]?([1-9]\d{1,14}[\,[0-9a-dA-D\*\#]+]*?)[^,]$`)

	subMatch := regex.FindStringSubmatch("+123")
	if len(subMatch) == 2 {
		fmt.Printf("test 1 %s\n", subMatch[1])
	}

	subMatch = regex.FindStringSubmatch("+12345")
	if len(subMatch) == 2 {
		fmt.Printf("test 2 %s\n", subMatch[1])
	}

	subMatch = regex.FindStringSubmatch("+12345,1334,2,34")
	if len(subMatch) == 2 {
		fmt.Printf("test 3 %s\n", subMatch[1])
	}

	subMatch = regex.FindStringSubmatch("12345,1334,2,34")
	if len(subMatch) == 2 {
		fmt.Printf("test 4 %s\n", subMatch[1])
	}

	subMatch = regex.FindStringSubmatch("0012345,1334,abcd")
	if len(subMatch) == 2 {
		fmt.Printf("test 5 %s\n", subMatch[1])
	}

	subMatch = regex.FindStringSubmatch("12345678,1334,abcdk")
	if len(subMatch) == 2 {
		fmt.Printf("test 6 %s\n", subMatch[1])
	}

	subMatch = regex.FindStringSubmatch("12345678112,1334,abcd,12345")
	if len(subMatch) == 2 {
		fmt.Printf("test 7 %s\n", subMatch[1])
	}

	subMatch = regex.FindStringSubmatch("12345678112,1334,abcd,")
	if len(subMatch) == 2 {
		fmt.Printf("test 8 %s\n", subMatch[1])
	}

	// subMatch = regex.FindSubmatch([]byte("++++++00000123456"))
	// if subMatch != nil {
	// 	fmt.Printf("%s\n", subMatch[1])
	// }

	// subMatch = regex.FindSubmatch([]byte("++++++00000123456string"))
	// if subMatch != nil {
	// 	fmt.Printf("%s\n", subMatch[1])
	// }

	// subMatch = regex.FindSubmatch([]byte("123456"))
	// if subMatch != nil {
	// 	fmt.Printf("%s\n", subMatch[1])
	// }

	// subMatch = regex.FindSubmatch([]byte("00123456"))
	// if subMatch != nil {
	// 	fmt.Printf("%s\n", subMatch[1])
	// }

	// subMatch = regex.FindSubmatch([]byte("+00123456"))
	// if subMatch != nil {
	// 	fmt.Printf("%s\n", subMatch[1])
	// }

	return
}

func normalizer(number string) string {
	// remove '+' or '00' prefix
	res := strings.TrimPrefix(number, "+")
	if res == number {
		res = strings.TrimPrefix(number, "00")
	}

	// if the number has any non-number characters, return the original
	_, err := strconv.Atoi(res)
	if err != nil {
		return number
	}

	return res
}
