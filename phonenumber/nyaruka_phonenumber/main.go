package main

import (
	"fmt"

	"github.com/nyaruka/phonenumbers"
)

func main() {
	// parse our phone number
	num, err := phonenumbers.Parse("+821121656521", "")
	if err != nil {
		fmt.Printf("Err: %v\n", err)
	}
	fmt.Printf("res: %v\n", num)

	// format it using national format
	formattedNum := phonenumbers.Format(num, phonenumbers.NATIONAL)
	fmt.Printf("res: %v\n", formattedNum)
}
