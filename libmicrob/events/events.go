package events

import (
	"fmt"
	"github.com/synw/microb/libmicrob/datatypes"
	"github.com/synw/microb/libmicrob/conf"
	"github.com/synw/microb/libmicrob/events/format"
)


var c = conf.GetConf()
var verbosity = c["verbosity"].(int)

func printMsg(event_class string, event *datatypes.Event) {
	msg := format.GetFormatedMsg(event_class, event)
	if event_class == "error" {
		msg = msg+" (from "+event.From+")"
	}
	fmt.Println(msg)
}

func Handle(event *datatypes.Event) {
	if verbosity > 0 {
		printMsg(event.Class, event)
	}
}

func PrintCommandReport(command *datatypes.Command) {
	fmt.Println(format.GetTime(), format.GetCommandReportMsg(command))
}

func Print(event_class string, from string, message string) {
	var d map[string]interface{}
	event := &datatypes.Event{event_class, from, message, d}
	printMsg(event_class, event)
}

func New(event_class string, from string, message string) {
	var d map[string]interface{}
	event := &datatypes.Event{event_class, from, message, d}
	Handle(event)
}

func ErrMsg(from string, msg string) {
	var d map[string]interface{}
	event := &datatypes.Event{"error", from, msg, d}
	Handle(event)
}

func State(from string, msg string) {
	New("state", from, msg)
}

func Debug(args ...string) {
	var msg string
	for _, arg := range(args) {
		msg = msg+arg+" "
	}
	New("debug", "runtime", msg)
}

func Error(from string, err error) {
	var d map[string]interface{}
	event := &datatypes.Event{"error", from, err.Error(), d}
	Handle(event)
}
