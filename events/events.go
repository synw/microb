package events

import (
	"fmt"
	"time"
	"strconv"
	"github.com/acmacalister/skittles"
	"github.com/synw/microb/datatypes"
	"github.com/synw/microb/conf"
)


var Config = conf.GetConf()
var Verbosity = Config["verbosity"].(int)

var info_m = " ["+skittles.Green("info")+"]"
var out_m = " ["+skittles.Yellow("out")+"]"
var command_m = " [=> "+skittles.Cyan("command")+"]"
var error_m = " ["+skittles.BoldRed("error")+"]"
var metric_m = " ["+skittles.Cyan("metric")+"]"

func getTime() string {
	t := time.Now()
	return t.Format("15:04:05")
}

func printMsg(event_class string, event *datatypes.Event) {
	t := getTime()
	out := t+" "
	msg := event.Message
	if event_class == "simple" {
		out = out+msg
	} else {
		if (event_class == "info") {
			out = out+info_m+" "+msg
		} else if (event_class == "out") {
			out = out+out_m+" "+msg
		} else if (event_class == "command") {
			out = out+command_m+" "+msg
		} else if (event_class == "error") {
			out = out+error_m+" "+msg
		} else if (event_class == "metric") {
			out = out+metric_m+" "+msg
		} else if (event_class == "request" || event_class == "request_error") {
			sc := event.Data["status_code"].(int)
			var sc_str string
			if sc == 404 {
				sc_str = skittles.Red(strconv.Itoa(sc))
			} else if sc == 200 {
				sc_str = skittles.Green(strconv.Itoa(sc))
			} else if sc == 500 {
				sc_str = skittles.BoldRed(strconv.Itoa(sc))
			}
			out = out+" "+sc_str+" "+msg
		}
	}
	fmt.Println(out)
}

func Handle(event *datatypes.Event) {
	event_class := event.Class
	if event.Class == "runtime_info" {
		event_class = "simple"
	}
	if event.Class == "request" {
		event_class = "request"
	}
	if Verbosity > 0 {
		printMsg(event_class, event)
	}
}

func New(event_class string, from string, message string) {
	var d map[string]interface{}
	event := &datatypes.Event{event_class, from, message, d}
	Handle(event)
}
