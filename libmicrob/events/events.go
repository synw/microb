package events

import (
	"fmt"
	"github.com/synw/microb/libmicrob/datatypes"
	"github.com/synw/microb/libmicrob/metadata"
	"github.com/synw/microb/libmicrob/events/format"
)


var Verbosity = metadata.GetVerbosity()

func printMsg(event_class string, event *datatypes.Event) {
	msg := format.GetFormatedMsg(event_class, event)
	if event_class == "error" {
		msg = msg+" (from "+event.From+")"
	}
	fmt.Println(msg)
}

func Handle(event *datatypes.Event) {
	if Verbosity > 0 {
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

func Error(from string, err error) {
	var d map[string]interface{}
	event := &datatypes.Event{"error", from, err.Error(), d}
	Handle(event)
}
