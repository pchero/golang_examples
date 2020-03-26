package main

import (
	"log/syslog"

	"github.com/sirupsen/logrus"
	lSyslog "github.com/sirupsen/logrus/hooks/syslog"
)

func main() {
	log := logrus.New()
	// log.Formatter = &logrus.TextFormatter{DisableTimestamp: true, DisableColors: true}
	log.SetFormatter(&logrus.TextFormatter{DisableTimestamp: true, DisableColors: true})
	// hook, err := lSyslog.NewSyslogHook("", "", syslog.LOG_INFO, "")
	// hook, err := lSyslog.NewSyslogHook("", "", syslog.LOG_DEBUG, "")
	hook, err := lSyslog.NewSyslogHook("", "", syslog.LOG_INFO, "test")
	log.SetLevel(logrus.DebugLevel)

	if err == nil {
		log.Hooks.Add(hook)
	}

	log.Error("Error: Hello world!")
	log.Debug("Debug: Hello world!")
}
