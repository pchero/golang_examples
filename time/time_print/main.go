package main

import (
	"fmt"
	"time"
)

func main() {

	timePrint(time.Now())

	// test := "2018-11-22 07:00:03"

	// from, _ := time.Parse(time.RFC3339, "2018-11-22T06:00:03Z")
	// until, _ := time.Parse(time.RFC3339, "2018-11-22T08:00:03Z")

	// layout := "2006-01-02 15:04:05"
	// t, _ := time.Parse(layout, test)
	// fmt.Println(t)

	// // check normal case
	// ret := timePrint(test, from, until)
	// fmt.Println(ret)

	// // check same tiem from
	// ret = timePrint("2018-11-22 06:00:03", from, until)
	// fmt.Println(ret)

	// // check same tiem until
	// ret = timePrint("2018-11-22 08:00:03", from, until)
	// fmt.Println(ret)

}

func timePrint(time time.Time) bool {
	if time.IsZero() == true {
		return false
	}

	fmt.Printf("%02d/%02d/%02d", time.Year(), time.Month(), time.Day())

	return true
}
