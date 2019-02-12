package main

import "fmt"

func main() {
	test := map[string]string{
		"Report0IAJitter":     "1024",
		"Report0FractionLost": "2048",
	}

	max := 1
	for i := 0; i < max; i++ {
		key := fmt.Sprintf("Report%dIAJitter", i)

		fmt.Printf("key: %s, val: %s\n", key, test[key])
	}
}
