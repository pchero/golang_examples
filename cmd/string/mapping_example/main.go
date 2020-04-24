package main

import "fmt"

func main() {
	t0 := map[string]string{}

	t0["test"] = "testVal"

	fmt.Println(t0["test"])
	fmt.Println(t0["testkey1"])

}
