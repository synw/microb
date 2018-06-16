package events

import (
	"errors"
	"fmt"
	"github.com/SKAhack/go-shortid"
	color "github.com/acmacalister/skittles"
	//"github.com/synw/microb/msgs"
	"github.com/synw/microb/types"
	"github.com/synw/terr"
	"time"
)

var g = shortid.Generator()

func State(service string, mesg string) *types.Event {
	event := build(service, mesg, nil, "state")
	handle(event)
	return event
}

func Info(service string, mesg string) *types.Event {
	event := build(service, mesg, nil, "info")
	handle(event)
	return event
}

func Warning(service string, mesg string, tr *terr.Trace) *types.Event {
	event := build(service, mesg, tr, "warning")
	handle(event)
	return event
}

func Error(service string, mesg string, tr *terr.Trace) *types.Event {
	event := build(service, mesg, tr, "error")
	handle(event)
	return event
}

func Fatal(service string, mesg string, tr *terr.Trace) *types.Event {
	event := build(service, mesg, tr, "fatal")
	handle(event)
	return event
}

func Panic(service string, mesg string, tr *terr.Trace) *types.Event {
	event := build(service, mesg, tr, "panic")
	handle(event)
	return event
}

func build(service string, mesg string, tr *terr.Trace, logLvls ...string) *types.Event {
	args := make(map[string]interface{})
	args["msg"] = mesg
	args["service"] = service
	args["trace"] = tr
	logLvl := "error"
	if len(logLvls) > 0 {
		logLvl = logLvls[0]
	}
	args["logLvl"] = logLvl
	args["class"] = logLvl
	event := new_(mesg, args)
	return event
}

func CmdExec(cmd *types.Cmd) {
	Cmd(cmd, true)
}

func Cmd(cmd *types.Cmd, out ...bool) {
	msg := color.BoldWhite(cmd.Name) + " from " + cmd.Service
	msg = msg + fmt.Sprintf("%s ", cmd.Date)
	cmd.LogMsg = cmd.Name + " received from service " + cmd.Service
	data := map[string]interface{}{
		"args":         cmd.Args,
		"returnValues": cmd.ReturnValues,
	}
	var tr *terr.Trace
	if cmd.ErrMsg != "" {
		err := errors.New(cmd.ErrMsg)
		tr = terr.New(err)
	}
	args := make(map[string]interface{})
	args["service"] = cmd.Service
	args["from"] = cmd.From
	//args["msg"] = msg
	args["trace"] = tr
	args["data"] = data
	args["cmd"] = cmd
	if len(out) > 0 {
		args["class"] = "command_out"
		var rvs string
		for _, val := range cmd.ReturnValues {
			rvs = rvs + val.(string)
			/*if i < (len(cmd.ReturnValues) + 1) {
				rvs = rvs + "\n"
			}*/
		}
		msg = rvs
		cmd.LogMsg = rvs
	} else {
		args["class"] = "command_in"
	}
	event := new_(msg, args)
	handle(event)
}

func new_(msg string, args ...map[string]interface{}) *types.Event {
	class, service, cmd, trace, logLvl, data := getEventArgs(args...)
	date := time.Now()
	id := g.Generate()
	event := &types.Event{id, class, date, msg, service, cmd, trace, logLvl, data}
	return event
}

func getEventArgs(args ...map[string]interface{}) (string, string, *types.Cmd, *terr.Trace, string, map[string]interface{}) {
	eclass := "default"
	var eservice string
	var ecmd *types.Cmd
	var etrace *terr.Trace
	logLvl := "default"
	edata := make(map[string]interface{})
	if len(args) > 0 {
		for _, arg := range args {
			for k, v := range arg {
				if k == "class" {
					eclass = v.(string)
				} else if k == "service" {
					eservice = v.(string)
				} else if k == "trace" {
					etrace = v.(*terr.Trace)
				} else if k == "cmd" {
					ecmd = v.(*types.Cmd)
				} else if k == "data" {
					edata = v.(map[string]interface{})
				} else if k == "logLvl" {
					logLvl = v.(string)
				}
			}
		}
	}
	return eclass, eservice, ecmd, etrace, logLvl, edata
}
