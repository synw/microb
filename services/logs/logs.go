package logs

import (
	"github.com/Sirupsen/logrus"
	"github.com/synw/microb/libmicrob/msgs"
	"github.com/synw/microb/libmicrob/types"
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

func Init(conf *types.Conf) {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	// initialize Redis connection
	tr := initRedis(conf)
	if tr != nil {
		tr.Fatal()
	}
	// initialize the lgos database
	tr = initDb(conf)
	if tr != nil {
		tr.Fatal()
	}
	// run the worker to process the logs
	key := "log_" + conf.Name
	go processLogs(key)
	//hook := NewCentHook(cli, logChans)
	hook := newRedisHook()
	logger.Hooks.Add(hook)
	msgs.Ok("Logger initialized")
	//logger.Out = ioutil.Discard
}

func Event(event *types.Event) {
	New(event.Service, event.LogLvl, event.Msg, event.Class)
}

func New(service string, level string, msg string, classes ...string) {
	class := "info"
	if len(classes) > 0 {
		class = classes[0]
	}
	now := time.Now().UnixNano() / int64(time.Millisecond)
	logobj := logger.WithFields(logrus.Fields{
		"service": service,
		"level":   level,
		"date":    now,
		"class":   class,
	})
	if level == "debug" {
		logobj.Debug(msg)
	} else if level == "info" {
		logobj.Info(msg)
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
