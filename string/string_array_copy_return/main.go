package main

import "fmt"

type stSimpleKeyValue struct {
	Key   string
	Value string
}

func main() {
	stTest1 := []stSimpleKeyValue{
		{Key: "key1", Value: "value1"},
		{Key: "key2", Value: "value2"},
		{Key: "key3", Value: "value3"},
	}

	tmp := returnArrItem(stTest1)
	fmt.Println(tmp)

	str1 := []string{
		"str1_key1",
		"str1_key2",
		"str1_key3",
		"str1_key4",
		"str1_key5",
	}

	fmt.Println(str1)
	return
}

func returnArrItem(stTests []stSimpleKeyValue) stSimpleKeyValue {
	var res []stSimpleKeyValue
	for _, stTest := range stTests {
		res = append(res, stTest)
	}

	return res[1]
}
