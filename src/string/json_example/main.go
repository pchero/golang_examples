package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	s, _ := json.Marshal(`{
		"TEST": 1
	}`)

	fmt.Println(string(s))

	t := map[string]string{
		"key1": "val1",
		"key2": "val2",
	}
	jTmp, _ := json.Marshal(t)
	fmt.Println(string(jTmp))
}
