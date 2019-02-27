package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	s := `{"test": "value"}`

	t, _ := json.Marshal(s)
	fmt.Println(string(t))

	// s =     byt := []byte(`{"num":6.13,"strs":["a","b"]}`)
	k := []byte(`{"num":6.1 ,"strs":["a","b"]}`)
	var tmpJSON map[string]interface{}
	err := json.Unmarshal(k, &tmpJSON)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(tmpJSON)
	fmt.Println(tmpJSON["num"])

	var objmap map[string]*json.RawMessage
	err = json.Unmarshal(t, &objmap)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(objmap)

	err = json.Unmarshal([]byte("TEST=1"), &tmpJSON)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(tmpJSON)
}
