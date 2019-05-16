// FtoC prints two Fahrenheit-to-Celsius conversions.
package main

import "fmt"

func main() {
	const freezingF, boilingF = 32.0, 212.0

	fmt.Printf("%gF = %gC\n", freezingF, fToC(freezingF)) // 32F = 0C
	fmt.Printf("%gF = %gC\n", boilingF, fToC(boilingF))   // 212F = 100C
}

func fToC(f float64) float64 {
	return (f - 32) * 5 / 9
}

// $ go run cmd/the_go_programming_language/chapter_2/ftoc/main.go
// 32F = 0C
// 212F = 100C
