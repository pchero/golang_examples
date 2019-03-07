package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

// TestStruct test struct
type TestStruct struct {
	Name string `json: name`
	Age  int    `json: age`
}

func main() {
	s := []string{"this", "is", "a", "joined", "string\n"}
	fmt.Println(strings.Join(s, " "))

	t := TestStruct{
		Name: "pchero",
		Age:  10,
	}

	fmt.Println(t)

	tmpJSON, err := json.Marshal(t)
	if err != nil {
		return
	}
	fmt.Println(string(tmpJSON))

}
