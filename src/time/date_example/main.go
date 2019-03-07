package main

import (
	"fmt"
	"time"
)

func main() {
	from, until := previousQuarterHour(time.Now().UTC())

	fmt.Printf("%s, %s\n", from, until)

}

func previousQuarterHour(t time.Time) (from, until time.Time) {
	if t.Minute() >= 45 {
		from = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 30, 0, 0, time.UTC)
		until = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 45, 0, 0, time.UTC)
	} else if t.Minute() >= 30 {
		from = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 15, 0, 0, time.UTC)
		until = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 30, 0, 0, time.UTC)
	} else if t.Minute() >= 15 {
		from = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, time.UTC)
		until = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 15, 0, 0, time.UTC)
	} else {
		from = time.Date(t.Year(), t.Month(), t.Day(), t.Hour()-1, 45, 0, 0, time.UTC)
		until = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, time.UTC)
	}

	return from, until
}
