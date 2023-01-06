package main

import "fmt"

func main() {
	m := map[string]string{
		"key1": "val1",
	}

	fmt.Println(m["key1"])
	fmt.Println(m["key2"])

	if m["key2"] == "" {
		fmt.Println("Wow")
	}
}
