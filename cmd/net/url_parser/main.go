package main

import "fmt"
import "net/url"

func main() {
	s := "postgres://user:pass@host.com:5432/path/test.wav?k=v#f"

	u, _ := url.Parse(s)
	fmt.Println(u.Path)

}