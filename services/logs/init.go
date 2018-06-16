package logs

import (
	"github.com/Sirupsen/logrus"
	"github.com/synw/microb/msgs"
	"github.com/synw/microb/types"
	"github.com/synw/terr"
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

func initService(dev bool, start bool) *terr.Trace {
	return nil
}

func Init(conf *types.Conf, state *types.State) *terr.Trace {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	// initialize the lgos database
	tr := initDb(conf)
	if tr != nil {
		tr.Fatal("Can not initilize logs database")
		return tr
	}
	// run the worker to process the logs
	key := "log_" + conf.Name
	go processLogs(key)
	channel := "$logs_" + conf.Name
	chook := NewCentHook(state.Cli, logChans, channel, state.Conf.Name)
	rhook := newRedisHook()
	logger.Hooks.Add(chook)
	msgs.State("Emiting logs on channel " + channel)
	logger.Hooks.Add(rhook)
	msgs.Ok("Logger initialized")
	logger.Out = ioutil.Discard
	return nil
}

func Event(event *types.Event) {
	if event.Cmd != nil {
		if event.Cmd.LogMsg != "" {
			event.Msg = event.Cmd.LogMsg
		}
	}
	New(event)
}

func New(event *types.Event) {
	if event.Class == "" {
		event.Class = "info"
	}
	now := time.Now().UnixNano() / int64(time.Millisecond)
	cmdName := ""
	cmdStatus := ""
	if event.Cmd != nil {
		cmdName = event.Cmd.Name
		cmdStatus = event.Cmd.Status
	}
	logobj := logger.WithFields(logrus.Fields{
		"service":        event.Service,
		"level":          event.LogLvl,
		"date":           now,
		"class":          event.Class,
		"command":        cmdName,
		"command_status": cmdStatus,
	})
	level := event.LogLvl
	msg := event.Msg
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
