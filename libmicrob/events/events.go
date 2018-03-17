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

func State(mesg string) *types.Event {
	args := make(map[string]interface{})
	args["class"] = "state"
	event := new_(mesg, args)
	return event
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
