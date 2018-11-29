package main

import (
	"fmt"
	"time"
)

func main() {
	test := "2018-11-22 07:00:03"

	from, _ := time.Parse(time.RFC3339, "2018-11-22T06:00:03Z")
	until, _ := time.Parse(time.RFC3339, "2018-11-22T08:00:03Z")

	layout := "2006-01-02 15:04:05"
	t, _ := time.Parse(layout, test)
	fmt.Println(t)

	// check normal case
	ret := timeCompare(test, from, until)
	fmt.Println(ret)

	// check same tiem from
	ret = timeCompare("2018-11-22 06:00:03", from, until)
	fmt.Println(ret)

	// check same tiem until
	ret = timeCompare("2018-11-22 08:00:03", from, until)
	fmt.Println(ret)

}

func timeCompare(timestamp string, from, until time.Time) bool {
	if len(timestamp) == 0 {
		return false
	}

	// parse
	layout := "2006-01-02 15:04:05"
	t, _ := time.Parse(layout, timestamp)

	return t.After(from) && t.Before(until)
}
