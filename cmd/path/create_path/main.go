package main

import (
	"fmt"
	"os"
	"path"
)

func main() {
	filename := "./test/test2/test.txt"
	err := os.MkdirAll(path.Dir(filename), 0755)
	if err != nil {
		fmt.Println(err)
	}

	file, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
}
