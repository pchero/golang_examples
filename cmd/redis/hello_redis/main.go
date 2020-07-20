package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/go-redis/redis"
)

var host = flag.String("host", "127.0.0.1:6379", "redis host address")

func main() {
	flag.Parse()
	ctx := context.Background()

	fmt.Printf("Hello world!\n")

	client := redis.NewClient(&redis.Options{
		Addr:     *host,
		Password: "",
		DB:       1,
	})
	// client.

	pong, err := client.Ping(ctx).Result()
	if err != nil {
		fmt.Printf("err: %v", err)
		return
	}
	fmt.Printf("result: %v\n", pong)
}
