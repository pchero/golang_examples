package main

import (
	"fmt"
	"regexp"
)

func main() {
	regex := regexp.MustCompile(`(?i)^sip:([^@\s\[\]:]+(:[^@\s\[\]:]+)?@)?[^@\s\[\]:]+(:[0-9]+)?$`)

	if ret := regex.Match([]byte("sip:test@test.com")); ret != true {
		fmt.Printf("Wrong\n")
	}

	if ret := regex.Match([]byte("sip:test@test.com:8080")); ret != true {
		fmt.Printf("Wrong1.\n")
	}

	if ret := regex.Match([]byte("sip:test@test.com :8080")); ret != true {
		fmt.Printf("Wrong2\n")
	}

	fmt.Printf("Good\n")
}
