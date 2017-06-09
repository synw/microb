package log

import (
	"github.com/Sirupsen/logrus"
	centhook "github.com/synw/logrus-centrifugo"
	"github.com/synw/microb/libmicrob/state"
	//"os"
	//"io"
	"time"
)

var logChans = map[string]string{
	"debug": "mysite_public",
	"info":  "mysite_public",
	"warn":  "mysite_public",
	"error": "mysite_public",
	"fatal": "mysite_public",
	"panic": "mysite_public",
}

var logger = logrus.New()

func Init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	var buf []byte
	var w = io.Writer(&buf)
	logrus.SetOutput(w)
	hook := centhook.New(state.Cli, logChans)
	logger.Hooks.Add(hook)
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
