package main

import (
	"encoding/json"
	"fmt"
)

func main() {

	msg := []byte(``)
	fmt.Println(msgValidate(msg))
	fmt.Println(msgValidate2(msg))

	msg = []byte(`{}`)
	fmt.Println(msgValidate(msg))
	fmt.Println(msgValidate2(msg))

	msg = []byte(`{"test1": "test2"}`)
	fmt.Println(msgValidate(msg))
	fmt.Println(msgValidate2(msg))

}

func msgValidate(msg []byte) bool {

	if len(msg) == 0 {
		return true
	}

	// we send only json formatted string.
	var js json.RawMessage
	if json.Unmarshal(msg, &js) != nil {
		return false
	}

	return true
}

func msgValidate2(msg []byte) bool {

	// zero sized message is ok.
	if len(msg) == 0 {
		return true
	}

	// we send only json formatted string.
	whocares := map[string]interface{}{}
	if err := json.Unmarshal(msg, &whocares); err != nil {
		return false
	}

	return true
}
