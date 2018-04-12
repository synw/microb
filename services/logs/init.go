package logs

import (
	"github.com/Sirupsen/logrus"
	"github.com/synw/microb/libmicrob/msgs"
	"github.com/synw/microb/libmicrob/types"
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

var Service *types.Service = &types.Service{
	"logs",
	getCmds(),
	initService,
}

var logger = logrus.New()

func initService(dev bool, start bool) error {
	// as it uses config the logger is initialized from the state
	return nil
}

func Init(conf *types.Conf) {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	// initialize the lgos database
	tr := initDb(conf)
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
	logger.Out = ioutil.Discard
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
	} else if level == "warning" {
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
