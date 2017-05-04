package events

import (
	"fmt"
	color "github.com/acmacalister/skittles"	
	"github.com/synw/terr"
	"github.com/synw/microb/libmicrob/datatypes"
	"github.com/synw/microb/libmicrob/state"
)

/*
func Error(trace *terr.Trace) {
	event := Err("error", "runtime", trace)
	handle(event)
}*/

// constructor

func Msg(class string, from string, msg string) *datatypes.Event {
	var data map[string]interface{}
	event := &datatypes.Event{class, from, msg, data, nil}
	handle(event)
	return event
}

func Err(class string, from string, tr ...*terr.Trace) *datatypes.Event {
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
	handle(event)
	return event
}

func Cmd(cmd *datatypes.Command) *datatypes.Event {
	msg := color.BoldWhite(cmd.Name)+" received "+fmt.Sprintf("%s", cmd.Date)
	if (cmd.Reason != "") {
		msg =msg+" ( "+cmd.Reason+" )"
	}
	args := make(map[string]interface{})
	args["args"] = cmd.Args
	event := &datatypes.Event{"command", cmd.From, msg, args, nil}
	handle(event)
	return event
}
/*
func Info(msg string) {
	var data map[string]interface{}
	event := &datatypes.Event{"info", "", msg, data, nil}
	handle(event)
}*/

// internal methods

func handle(event *datatypes.Event) {
	if state.Verbosity > 0 {
		fmt.Println(getMsg(event))
	}
}

func getMsg(event *datatypes.Event) string {
	out := getFormatedMsgNoTime(event)
	return out
}

func getEventOutputFlags() map[string]string {
	output_flags := make(map[string]string)
	output_flags["info"] = "["+color.Blue("info")+"]"
	output_flags["event"] = "["+color.Yellow("event")+"]"
	output_flags["command"] = "[=> "+color.Cyan("command")+"]"
	output_flags["error"] = ""
	output_flags["metric"] = "["+color.Cyan("metric")+"]"
	output_flags["state"] = "["+color.Yellow("state")+"]"
	output_flags["debug"] = "["+color.Magenta("debug")+"]"
	output_flags["request"] = ""
	output_flags["request_error"] = ""
	output_flags["runtime_info"] = ""
	return output_flags
}

func getFormatedMsgNoTime(event *datatypes.Event) string {
	output_flags := getEventOutputFlags()
	label := output_flags[event.Class]
	sep := " "
	if label == "" {
		sep = ""
	}
	msg := event.Message
	out := label+sep+msg
	return out
}
