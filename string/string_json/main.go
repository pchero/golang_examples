package main

import (
	"encoding/json"
	"fmt"
)

type response1 struct {
	test1 *response2 `json: "-"`
	Key1  string     `json: "key1"`
}

type response2 struct {
	Key1 string `json: "key1"`
	Key2 string `json: "key"`
}

func main() {
	test1 := &response2{
		Key1: "val1",
		Key2: "val2",
	}

	s, err := json.Marshal(test1)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(s))

	var str map[string]interface{}
	err = json.Unmarshal(s, &str)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(str)

	test2 := &response1{
		test1: test1,
		Key1:  "test",
	}

	s, err = json.Marshal(test2)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(s))

}
