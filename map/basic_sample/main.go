package main

import "fmt"

func main() {
	test := map[string]string{
		"key1": "val1",
	}

	fmt.Printf("tmp: %v\n", test)

	val, ok := test[""]
	fmt.Printf("empty val: %v, ok: %v\n", val, ok)

	test[""] = "val2"
	val, ok = test[""]
	fmt.Printf("empty val: %v, ok: %v\n", val, ok)
}
