package events

import (
	"fmt"
	"time"
	"github.com/acmacalister/skittles"
	"github.com/synw/microb/datatypes"
)


var info_m = " ["+skittles.Green("info")+"]"
var out_m = " ["+skittles.Yellow("out")+"]"
var command_m = " [=> "+skittles.Cyan("command")+"]"
var error_m = " ["+skittles.BoldRed("error")+"]"
var metric_m = " ["+skittles.Cyan("metric")+"]"

func getTime() string {
	t := time.Now()
	return t.Format("15:04:05")
}

func printMsg(print_type string, msg string) {
	t := getTime()
	out := t+" "
	if print_type == "simple" {
		out = out+msg
	} else {
		if (print_type == "info") {
			out = out+info_m+" "+msg
		} else if (print_type == "out") {
			out = out+out_m+" "+msg
		} else if (print_type == "command") {
			out = out+command_m+" "+msg
		} else if (print_type == "error") {
			out = out+error_m+" "+msg
		} else if (print_type == "metric") {
			out = out+metric_m+" "+msg
		}
	}
	fmt.Println(out)
}

func NewEvent(event_class string, from string, message string) *datatypes.Event {
	var d map[string]interface{}
	event := &datatypes.Event{event_class, from, message, d}
	return event
}

func Handle(event *datatypes.Event, verbosity int) {
	print_type := event.Class
	if event.Class == "runtime_info" {
		print_type = "simple"
	}
	if verbosity > 0 {
		printMsg(print_type, event.Message)
	}
}
