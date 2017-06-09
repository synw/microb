package events

import (
	"errors"
	"fmt"
	color "github.com/acmacalister/skittles"
	"github.com/synw/microb/libmicrob/datatypes"
	"github.com/synw/microb/libmicrob/log"
	"github.com/synw/microb/libmicrob/state"
	"time"
)

func New(class string, service string, from string, msg string, err error, data ...map[string]interface{}) *datatypes.Event {
	dataset := make(map[string]interface{})
	if len(data) > 0 {
		dataset = data[0]
	}
	now := time.Now()
	event := &datatypes.Event{class, from, service, now, msg, err, dataset}
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

func Cmd(cmd *datatypes.Command) {
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

// internal methods

func handle(event *datatypes.Event) {
	if state.Verbosity > 0 {
		fmt.Println(getMsg(event))
	}
	log.New(event.Service, state.Conf["name"].(string), event.Class, event.Msg)
}

func getMsg(event *datatypes.Event) string {
	out := getFormatedMsgNoTime(event)
	return out
}

func getFormatedMsgNoTime(event *datatypes.Event) string {
	output_flags := getEventOutputFlags()
	label := output_flags[event.Class]
	sep := " "
	if label == "" {
		sep = ""
	}
	msg := event.Msg
	out := label + sep + msg
	return out
}

func getEventOutputFlags() map[string]string {
	output_flags := make(map[string]string)
	output_flags["ready"] = "[" + color.Green("ready") + "]"
	output_flags["info"] = "[" + color.Blue("info") + "]"
	output_flags["event"] = "[" + color.Yellow("event") + "]"
	output_flags["command"] = "[=> " + color.Cyan("command") + "]"
	output_flags["error"] = ""
	output_flags["metric"] = "[" + color.Cyan("metric") + "]"
	output_flags["state"] = "[" + color.Yellow("state") + "]"
	output_flags["debug"] = "[" + color.Magenta("debug") + "]"
	return output_flags
}
