// https://gobyexample.com/time-formatting-parsing

// Go supports time formatting and parsing via
// pattern-based layouts.

package main

import (
	"fmt"
	"time"
)

func main() {
	p := fmt.Println

	// Here's basic example of formatting a time
	// according to RFC3339, using the corresponding layout
	// constant.
	t := time.Now()
	p(t.Format(time.RFC3339))

	// Time parsing uses the same layout values as 'Format'.
	t1, e := time.Parse(
		time.RFC3339,
		"2012-11-01T22:08:41+00:00",
	)
	if e != nil {
		p(e)
	}
	p(t1)

	// 'Format' and 'Parse' use example-based layouts. Usually
	// you'll use a constant from 'time' for these layouts, but
	// you can also supply custom layouts. Layouts must use the
	// reference time 'Mon Jan 2 15:04:05 MST 2006' to show the
	// patterns with which to format/parse a given time/string.
	// The example time must be exactly as shown: the year 2006,
	// 15 for the hour, Monday for the day of the week, etc.
	p(t.Format("3:04PM"))
	p(t.Format("Mon Jan _2 15:04:05 2006"))
	p(t.Format("2006-01-02T15:04:05.999999-07:00"))
	form := "3 04 PM"
	t2, e := time.Parse(form, "8 41 PM")
	p(t2)

	// For purely numeric representations you can also
	// use standard string formatting with the extracted
	// components of the time value.
	fmt.Printf("%d-%02d-%02dT%02d:%02d:%02d-00:00\n",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second(),
	)

	// 'Parse' will return an error on malformed input
	// explaining the parsing problem.
	ansic := "Mon Jan _2 15:04:05 2006"
	_, e = time.Parse(ansic, "8:14 PM")
	p(e)
}

// $ go run cmd/gobyexample_2/time_formatting_parsing/main.go
// 2019-04-02T09:11:16+02:00
// 2012-11-01 22:08:41 +0000 +0000
// 9:11AM
// Tue Apr  2 09:11:16 2019
// 2019-04-02T09:11:16.539296+02:00
// 0000-01-01 20:41:00 +0000 UTC
// 2019-04-02T09:11:16-00:00
// parsing time "8:14 PM" as "Mon Jan _2 15:04:05 2006": cannot parse "8:14 PM" as "Mon"
