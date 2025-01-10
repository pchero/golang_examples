package main

import (
	"fmt"
	"time"

	"github.com/BorisBorshevsky/timemock"
)

func main() {
	fmt.Println("timemock.Now():", timemock.Now())
	fmt.Println("timemock.Now():", time.Now())

}
