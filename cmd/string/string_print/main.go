package main

import "fmt"

func main() {
	fmt.Printf("Hello world!\n")

	s := "Hello world!\n"[0:3]
	fmt.Println(s) // Hel

	d := s[:len(s)-1]
	fmt.Print(d) // He
}
