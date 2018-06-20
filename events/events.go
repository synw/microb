package events

import (
	"fmt"
	"github.com/SKAhack/go-shortid"
	color "github.com/acmacalister/skittles"
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

/* Commands */

func CmdIn(cmd *types.Cmd) {
	msg := color.BoldWhite(cmd.Name) + " from " + cmd.Service
	msg = msg + fmt.Sprintf("%s ", cmd.Date)
	data := map[string]interface{}{
		"args":         cmd.Args,
		"returnValues": cmd.ReturnValues,
	}
	args := make(map[string]interface{})
	args["service"] = cmd.Service
	args["from"] = cmd.From
	args["data"] = data
	args["cmd"] = cmd
	args["class"] = "command_in"
	event := build_(msg, args)
	handle(event)
}

func CmdOut(cmd *types.Cmd) {
	msg := color.BoldWhite(cmd.Name) + " from " + cmd.Service
	msg = msg + fmt.Sprintf("%s ", cmd.Date)
	data := map[string]interface{}{
		"args":         cmd.Args,
		"returnValues": cmd.ReturnValues,
	}
	args := make(map[string]interface{})
	args["service"] = cmd.Service
	args["from"] = cmd.From
	if cmd.Trace != nil {
		args["trace"] = cmd.Trace
	}
	args["data"] = data
	args["cmd"] = cmd
	args["class"] = "command_out"
	var rvs string
	for _, val := range cmd.ReturnValues {
		rvs = rvs + val.(string)
		msg = rvs
	}
	event := build_(msg, args)
	handle(event)
}

func build(service string, mesg string, tr *terr.Trace, logLvls ...string) *types.Event {
	args := make(map[string]interface{})
	args["msg"] = mesg
	args["service"] = service
	args["trace"] = tr
	//args.["cmd"] = cmd
	logLvl := "error"
	if len(logLvls) > 0 {
		logLvl = logLvls[0]
	}
	args["logLvl"] = logLvl
	args["class"] = logLvl
	event := build_(mesg, args)
	return event
}

func build_(msg string, args ...map[string]interface{}) *types.Event {
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
				} else if k == "cmd" {
					ecmd = v.(*types.Cmd)
				} else if k == "trace" {
					etrace = v.(*terr.Trace)
				} else if k == "logLvl" {
					logLvl = v.(string)
				} else if k == "data" {
					edata = v.(map[string]interface{})
				}
			}
		}
	}
	return eclass, eservice, ecmd, etrace, logLvl, edata
}
