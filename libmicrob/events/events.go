package events

import (
	"errors"
	"fmt"
	"github.com/SKAhack/go-shortid"
	color "github.com/acmacalister/skittles"
	"github.com/synw/microb/libmicrob/state"
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
	event := &types.Event{id, class, from, service, now, msg, err, dataset}
	handle(event)
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
	err := tr.ToErr()
	_ = New("error", service, from, msg, err, dataset)
}

func State(service string, from string, msg string, err error, data ...map[string]interface{}) {
	dataset := make(map[string]interface{})
	if len(data) > 0 {
		dataset = data[0]
	}
	_ = New("state", service, from, msg, err, dataset)
}

func Ready(service string, from string, msg string, err error, data ...map[string]interface{}) {
	dataset := make(map[string]interface{})
	if len(data) > 0 {
		dataset = data[0]
	}
	_ = New("ready", service, from, msg, err, dataset)
}

func Cmd(cmd *types.Command) {
	msg := color.BoldWhite(cmd.Name) + " received " + fmt.Sprintf("%s", cmd.Date)
	if cmd.Reason != "" {
		msg = msg + " ( " + cmd.Reason + " )"
	}
	data := map[string]interface{}{
		"args":         cmd.Args,
		"returnValues": cmd.ReturnValues,
	}
	var err error
	if cmd.ErrMsg != "" {
		err = errors.New(cmd.ErrMsg)
	}
	_ = New("command", cmd.Service, cmd.From, msg, err, data)
}

func CmdExec(cmd *types.Command) {
	data := map[string]interface{}{
		"args":         cmd.Args,
		"returnValues": cmd.ReturnValues,
	}
	var err error
	if cmd.ErrMsg != "" {
		err = errors.New(cmd.ErrMsg)
	}
	status := cmd.Status
	if status == "error" {
		status = color.BoldRed("error")
		if state.Verbosity > 0 {
			fmt.Println(" ->", status, cmd.Trace.Format())
		}
	} else if status == "success" {
		status = color.Green("success")
		if state.Verbosity > 0 {
			fmt.Println(" ->", status, cmd.ReturnValues)
		}
	}
	msg := ""
	_ = New("commandExec", cmd.Service, cmd.From, msg, err, data)
}
