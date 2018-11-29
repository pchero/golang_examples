package main

// package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"
)

type cdrAsterisk struct {
	uniqueid string
	started  string
	ended    string
	duration int
}

func (c cdrAsterisk) isEmpty() bool {
	return c.uniqueid == ""
}

type record struct {
	total    int
	duration int
}

func main() {

	// test
	f, err := os.Open("./test_1.cdr")
	if err != nil {
		println(err)
	}
	defer f.Close()

	res, err := parse(bufio.NewReader(f))
	fmt.Printf("total[%d], duration[%d]\n", res.total, res.duration)

	// test
	f, err = os.Open("./voice-asterisk-production-europe-west1-d-10%2Fcdr%2F2018%2F11%2F22%2Fcdr-20181122080001.csv")
	if err != nil {
		println(err)
	}
	defer f.Close()

	res, err = parse(bufio.NewReader(f))
	fmt.Printf("total[%d], duration[%d]\n", res.total, res.duration)

}

func parseAsterisk(reader *bufio.Reader, from, until time.Time) (record, error) {
	var record record

	cdrs := make(map[string]cdrAsterisk)
	csvReader := csv.NewReader(reader)

	// create cdr list
	for {
		line, err := csvReader.Read()
		if err != nil {
			break
		}

		// leg_id, clid, src, dst, dcontext, channel, dstchannel, lastapp, lastdata, start, answer, end, duration, billsec, disposition, amaflag, accountcode, uniqueid, userfiled, sequence
		cdr := cdrAsterisk{
			uniqueid: line[17],
			started:  line[9],
			ended:    line[11],
		}
		if cdr.duration, err = strconv.Atoi(line[12]); err != nil {
			continue
		}

		tmpCdr := cdrs[cdr.uniqueid]
		if tmpCdr.isEmpty() {
			// if cdr is not exist,
			// just insert it.
			cdrs[cdr.uniqueid] = cdr
			continue
		}

		// if cdr is exist,
		// it's the same call. add the call duration.
		tmpCdr.duration += cdr.duration
		tmpCdr.ended = cdr.ended
	}

	// count cdr and sum the duration
	for _, tmpCdr := range cdrs {

		// check is valid call
		if isValidCallAsterisk(tmpCdr, from, until) == false {
			continue
		}

		// increase count and duration
		record.duration += tmpCdr.duration
		record.total++
	}

	return record, nil
}

func isValidCallAsterisk(cdr cdrAsterisk, from, until time.Time) bool {
	// call has to be successful
	if cdr.duration == 0 {
		return false
	}

	// asterisk's timestmp layout
	// "2006-01-02 15:04:05"
	started, _ := time.Parse("2006-01-02 15:04:05", cdr.started)
	ended, _ := time.Parse("2006-01-02 15:04:05", cdr.ended)

	// started and end date need to be within our date range
	if started.After(from) != true || ended.Before(until) != true {
		return false
	}

	return true
}
