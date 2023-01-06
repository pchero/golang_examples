package main

import (
	"fmt"
	"time"
)

type test struct {
	duration int
}

func (t *test) Sleep() {
	time.Sleep(time.Second * 3)
}

func main() {
	fmt.Printf("hello world\n")

	t := test{
		duration: 10,
	}

	go func(t *test) {
		t.Sleep()
		fmt.Printf("1 done\n")
	}(&t)

	go func(t *test) {
		t.Sleep()
		fmt.Printf("2 done\n")
	}(&t)

	time.Sleep(time.Second * 10)
}
