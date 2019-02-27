package main

import (
	"encoding/json"
	"fmt"
)

type tmpStruct struct {
	key1 string
	key2 string
}

func main() {
	tmpTest := tmpStruct{
		key1: "test1",
		key2: "test2",
	}

	fmt.Println(tmpTest)

	tmp, err := json.Marshal(tmpTest)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(tmp)

	var tmpTest2 
	tmpTest2, err := json.Unmarshal(tmp)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(tmpTest2)
}
