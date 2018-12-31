package main

import "fmt"

func main() {
	data := make(map[int]int)

	for i := 0; i < 5; i++ {
		data[i] = i
	}

	// print values
	for i := 0; i < len(data); i++ {
		fmt.Print(data[i])
	}
	fmt.Println()

	// update values(wrong)
	for i := 0; i < len(data); i++ {
		tmp := data[i]
		tmp++
	}

	// print values
	for i := 0; i < len(data); i++ {
		fmt.Print(data[i])
	}
	fmt.Println()

	// update values(correct)
	for i, _ := range data {
		data[i] = data[i] + 1
	}

	// print values
	for i, _ := range data {
		fmt.Print(data[i])
	}
	fmt.Println()

}
