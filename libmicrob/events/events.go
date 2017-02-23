package events

import (
	"fmt"
	"github.com/synw/microb/libmicrob/datatypes"
	"github.com/synw/microb/libmicrob/state"
	"github.com/synw/microb/libmicrob/events/format"
)


func Handle(event *datatypes.Event) {
	if state.Verbosity > 0 {
		printMsg(event.Class, event)
	}
}

func PrintCommandReport(command *datatypes.Command) {
	if state.Verbosity > 0 {
		fmt.Println(format.GetTime(), format.GetCommandReportMsg(command))
	}
}

func Print(event_class string, from string, message string) {
	var d map[string]interface{}
	event := &datatypes.Event{event_class, from, message, d}
	if state.Verbosity > 0 {
		printMsg(event_class, event)
	}
}

func New(event_class string, from string, message string) {
	var d map[string]interface{}
	event := &datatypes.Event{event_class, from, message, d}
	Handle(event)
}

func State(from string, msg string) {
	New("state", from, msg)
}

func Debug(args ...string) {
	var msg string
	for i, arg := range(args) {
		if ((i+1) < len(args)) && (i>1) {
			msg = " "+msg
		}
		msg = msg+arg
	}
	New("debug", "runtime", msg)
}

func Error(from string, err error) {
	var d map[string]interface{}
	event := &datatypes.Event{"error", from, err.Error(), d}
	Handle(event)
}

func ErrMsg(from string, msg string) {
	var d map[string]interface{}
	event := &datatypes.Event{"error", from, msg, d}
	Handle(event)
}

func Err(from string, msg string, err error) {
	var d map[string]interface{}
	fmsg := msg+": "+err.Error()
	event := &datatypes.Event{"error", from, fmsg, d}
	Handle(event)
}

func printMsg(event_class string, event *datatypes.Event) {
	msg := format.GetFormatedMsg(event_class, event)
	if event_class == "error" {
		msg = msg+" (from "+event.From+")"
	}
	fmt.Println(msg)
}
