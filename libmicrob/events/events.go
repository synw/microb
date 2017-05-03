package events

import (
	"fmt"
	"errors"
	color "github.com/acmacalister/skittles"	
	"github.com/synw/terr"
	"github.com/synw/microb/libmicrob/datatypes"
	"github.com/synw/microb/libmicrob/state"
)


func Handle(event *datatypes.Event) {
	if state.Verbosity > 0 {
		fmt.Println(getMsg(event))
	}
}

// event classes

func Error(trace *terr.Trace) {
	event := New("error", "runtime", trace)
	Handle(event)
}

// constructor

func Msg(class string, from string, msg string) *datatypes.Event {
	err := errors.New(msg)
	tr := terr.New(from, err)
	event := New(class, from, tr)
	return event
}

func New(class string, from string, tr ...*terr.Trace) *datatypes.Event {
	var data map[string]interface{}
	var trace *terr.Trace
	var msg string
	if class == "error" {
		if len(tr) > 0 {
			trace = tr[0]
			msg = trace.Formatc()
		}
	}
	event := &datatypes.Event{class, from, msg, data, trace}
	return event
}

// internal methods

func getMsg(event *datatypes.Event) string {
	out := getFormatedMsgNoTime(event)
	return out
}

func getEventOutputFlags() map[string]string {
	output_flags := make(map[string]string)
	output_flags["info"] = "["+color.Green("info")+"]"
	output_flags["event"] = "["+color.Yellow("event")+"]"
	output_flags["command"] = "[=> "+color.Cyan("command")+"]"
	output_flags["error"] = ""
	output_flags["metric"] = "["+color.Cyan("metric")+"]"
	output_flags["state"] = "["+color.Yellow("state")+"]"
	output_flags["debug"] = "["+color.Yellow("debug")+"]"
	output_flags["request"] = ""
	output_flags["request_error"] = ""
	output_flags["runtime_info"] = ""
	return output_flags
}

func getFormatedMsgNoTime(event *datatypes.Event) string {
	var out string
	msg := event.Message
	output_flags := getEventOutputFlags()
	label := output_flags[event.Class]
	sep := " "
	if label == "" {
		sep = ""
	}
	out = out+sep+msg
	return out
}
