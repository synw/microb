package events

import (
	"errors"
	"fmt"
	"github.com/SKAhack/go-shortid"
	color "github.com/acmacalister/skittles"
	"github.com/synw/microb/libmicrob/types"
	"github.com/synw/terr"
	"time"
)

var g = shortid.Generator()

func New(class string, service string, from string, msg string, err error, data ...map[string]interface{}) *types.Event {
	dataset := make(map[string]interface{})
	if len(data) > 0 {
		dataset = data[0]
	}
	now := time.Now()
	id := g.Generate()
	var cmd *types.Cmd
	event := &types.Event{id, class, now, msg, service, cmd, nil, dataset}
	handle(event)
	return event
}

func State(mesg string) *types.Event {
	args := make(map[string]interface{})
	args["class"] = "state"
	event := new_(mesg, args)
	return event
}

func Err(service string, from string, msg string, err error, data ...map[string]interface{}) {
	dataset := make(map[string]interface{})
	if len(data) > 0 {
		dataset = data[0]
	}
	_ = New("error", service, from, msg, err, dataset)
}

func Terr(service string, from string, msg string, tr *terr.Trace, data ...map[string]interface{}) {
	dataset := make(map[string]interface{})
	if len(data) > 0 {
		dataset = data[0]
	}
	err := errors.New(tr.Formatc())
	_ = New("error", service, from, msg, err, dataset)
}

func CmdExec(cmd *types.Cmd) {
	Cmd(cmd, true)
}

func Cmd(cmd *types.Cmd, out ...bool) {
	msg := color.BoldWhite(cmd.Name) + " received " + fmt.Sprintf("%s", cmd.Date)
	data := map[string]interface{}{
		"args":         cmd.Args,
		"returnValues": cmd.ReturnValues,
	}
	var err error
	if cmd.ErrMsg != "" {
		err = errors.New(cmd.ErrMsg)
	}
	args := make(map[string]interface{})
	args["service"] = cmd.Service
	args["from"] = cmd.From
	args["msg"] = msg
	args["err"] = err
	args["data"] = data
	args["cmd"] = cmd
	if len(out) > 0 {
		args["class"] = "command_out"
	} else {
		args["class"] = "command_in"
	}
	_ = new_(msg, args)
}

func new_(msg string, args ...map[string]interface{}) *types.Event {
	class, service, cmd, trace, data := getEventArgs(args...)
	date := time.Now()
	id := g.Generate()
	event := &types.Event{id, class, date, msg, service, cmd, trace, data}
	handle(event)
	return event
}

func getEventArgs(args ...map[string]interface{}) (string, string, *types.Cmd, *terr.Trace, map[string]interface{}) {
	eclass := "default"
	var eservice string
	var ecmd *types.Cmd
	var etrace *terr.Trace
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
				}
			}
		}
	}
	return eclass, eservice, ecmd, etrace, edata
}
