package main

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

// func init() {
// 	// log as json instead of the default ascii formatter
// 	github.com
// }

// func main() {
// var log = logrus.New()

// log.Formatter = new(logrus.JSONFormatter)
// log.Formatter = new(logrus.TextFormatter)
// log.Level = logrus.TraceLevel
// log.Out = os.Stdout

// 	log.WithField(log.Fields{
// 		"animal": "walrus",
// 	}).Info("A walrus appears")
// }

func main() {
	log.SetFormatter(&log.JSONFormatter{})

	// tmp := []byte("testtttt      \n\000\012")
	tmp := []byte("testtttt")

	log.WithFields(log.Fields{
		"animal": "walrus",
		"test":   tmp,
	}).Info("A walrus appears")

	fmt.Println(tmp)
}
