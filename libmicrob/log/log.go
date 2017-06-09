package log

import (
	"github.com/Sirupsen/logrus"
	"github.com/synw/microb/libmicrob/state"
	"io/ioutil"
	"time"
)

var logChans = map[string]string{
	"debug": "$logchan",
	"info":  "$logchan",
	"warn":  "$logchan",
	"error": "$logchan",
	"fatal": "$logchan",
	"panic": "$logchan",
}

var logger = logrus.New()

func Init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	hook := NewHook(state.Cli, logChans)
	logger.Hooks.Add(hook)
	logger.Out = ioutil.Discard
}

func New(service string, instance string, level string, msg string) {
	now := time.Now().UnixNano() / int64(time.Millisecond)
	logobj := logger.WithFields(logrus.Fields{
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
