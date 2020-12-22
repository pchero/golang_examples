package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	fmt.Printf("Hello world\n")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	chanRes := make(chan *string)
	chanStop := make(chan bool)

	// do work
	go func() {
		for {
			select {
			case <-chanStop:
				return

			default:
				deadline, ok := ctx.Deadline()
				fmt.Printf("working... deadline: %v, ok: %v, \n", deadline, ok)
				time.Sleep(time.Millisecond * 100)
				continue
			}
		}
	}()

	select {
	case res := <-chanRes:
		fmt.Printf("Res: %s", *res)

	case <-ctx.Done():
		chanStop <- true
		fmt.Printf("Finished by timeout. err: %v\n", ctx.Err())
	}

	fmt.Printf("Wait another 3 sec..\n")
	time.Sleep(time.Second * 3)
}
