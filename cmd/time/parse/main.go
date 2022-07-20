package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	// convert time string to time.Time type

	timeStampString := "Dec 29, 2014 at 7:54pm (SGT)"
	layOut := "Jan 2, 2006 at 3:04pm (MST)"
	timeStamp, err := time.Parse(layOut, timeStampString)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	hr, min, sec := timeStamp.Clock()
	fmt.Printf("Clock : [%d]hour : [%d]minutes : [%d] seconds \n", hr, min, sec)

	year, month, day := timeStamp.Date()
	fmt.Printf("Date : [%d]year : [%d]month : [%d]day \n", year, month, day)

	now := time.Now().UTC().String()
	res := strings.TrimSuffix(now, " +0000 UTC")
	fmt.Printf("res: %s\n", res)

	test := "2021-04-18 03:22:17.994000"
	testLayOut := "2006-01-02 15:05:05.000000"
	testTimestamp, _ := time.Parse(testLayOut, test)
	fmt.Printf("Res: %v\n", testTimestamp)

	// now := time.Now()
	timeNow := time.Now()
	duration := timeNow.Sub(testTimestamp)
	fmt.Printf("Duration: %v\n", duration)
	// layOut :=
}
