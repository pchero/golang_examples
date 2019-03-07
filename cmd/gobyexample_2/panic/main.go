package main

import (
	"fmt"
	"os"
)

func main() {
	panic("a problem")

	_, err := os.Create("/tmp/file")
	if err != nil {
		panic(err)
	}

	fmt.Println("test gogo")
}
