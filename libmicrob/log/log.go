package log

import (
	log "github.com/Sirupsen/logrus"
	"os"
	"time"
)

func InitLogger() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	//log.SetLevel(log.DebugLevel)
}

func New(service string, instance string, level string, msg string) {
	now := time.Now().UnixNano() / int64(time.Millisecond)
	logobj := log.WithFields(log.Fields{
		"service":  service,
		"instance": instance,
		"level":    level,
		"date":     now,
	})
	if level == "debug" {
		logobj.Debug(msg)
	} else if level == "warn" {
		logobj.Warn(msg)
	} else if level == "error" {
		logobj.Error(msg)
	} else if level == "fatal" {
		logobj.Fatal(msg)
	} else if level == "panic" {
		logobj.Panic(msg)
	} else {
		logobj.Info(msg)
	}
}
