package main

import (
	"fmt"
	"regexp"
	"time"

	"github.com/sirupsen/logrus"
)

var re = regexp.MustCompile(`(\d{4})(\d{2})(\d{2})(\d{2})(\d{2})`)
var matchFile = regexp.MustCompile
var timeOffset = "+00:00" // Hardcoded UTC offset

func main() {
	parseString("Buckets/voice-bucket-production-europe-west1-d/voice-asterisk-production-europe-west1-d-10/cdr/2018/12/05/cdr-20181205120001.csv")

	parseString("Buckets/voice-bucket-production-europe-west1-d/voice-kamailio-production-europe-west1-d-01/cdr/2018/12/05/acc_cdrs_1.log.201812051215.tar.gz.enc")
}

func parseString(str string) {

	from := time.Unix(1544009440, 0)

	matchString := fmt.Sprintf("(.*)/cdr/%02d/%02d/%02d/(.*)", from.Year(), from.Month(), from.Day())
	fmt.Println(matchString)
	r, _ := regexp.Compile(matchString)
	match := r.MatchString(str)
	if match != true {
		return
	}

	fz := re.FindStringSubmatch(str)

	if len(fz) != 6 {
		return
	}

	fd, err := time.Parse(time.RFC3339, fmt.Sprintf("%s-%s-%sT%s:%s:00%s", fz[1], fz[2], fz[3], fz[4], fz[5], timeOffset))
	if err != nil {
		logrus.Warnf("Can not parse date %s: %s\n", fz, err)
		return
	}

	fmt.Println(fd.Hour(), from.Hour())
	if fd.Hour() < from.Hour() {
		return
	}

	fmt.Println(fd)
}
