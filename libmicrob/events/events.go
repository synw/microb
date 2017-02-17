package events

import (
	"fmt"
	"time"
	"strconv"
	"github.com/acmacalister/skittles"
	"github.com/synw/microb/libmicrob/datatypes"
	"github.com/synw/microb/libmicrob/metadata"
)


var Verbosity = metadata.GetVerbosity()

func getEventOutputFlags() map[string]string {
	output_flags := make(map[string]string)
	output_flags["info"] = "["+skittles.Cyan("info")+"]"
	output_flags["event"] = "["+skittles.Yellow("event")+"]"
	output_flags["command"] = "[=> "+skittles.Cyan("command")+"]"
	output_flags["error"] = "["+skittles.BoldRed("error")+"]"
	output_flags["metric"] = "["+skittles.Cyan("metric")+"]"
	output_flags["request"] = ""
	output_flags["request_error"] = ""
	output_flags["runtime_info"] = ""
	return output_flags
}

func getTime() string {
	t := time.Now()
	return t.Format("15:04:05")
}

func formatStatusCode(sc int) string {
	var sc_str string
	if sc == 404 {
		sc_str = skittles.Red(strconv.Itoa(sc))
	} else if sc == 200 {
		sc_str = skittles.Green(strconv.Itoa(sc))
	} else if sc == 500 {
		sc_str = skittles.BoldRed(strconv.Itoa(sc))
	}
	return sc_str
}

func printMsg(event_class string, event *datatypes.Event) {
	t := getTime()
	out := t+" "
	msg := event.Message
	output_flags := getEventOutputFlags()
	out = out+output_flags[event_class]
	if (event_class == "request" || event_class == "request_error") {
		sc := event.Data["status_code"].(int)
		out = out+formatStatusCode(sc)
	}
	out = out+" "+msg
	fmt.Println(out)
}

func Handle(event *datatypes.Event) {
	if Verbosity > 0 {
		printMsg(event.Class, event)
	}
}

func GetCommandReportMsg(command *datatypes.Command) string {
	var status string
	var vals string
	msg := ""
	if command.Status == "success" {
		status = skittles.Green("ok")	
	} else if command.Status == "error" {
		status = skittles.Red("error")
	} else {
		status = command.Status
	}
	if len(command.ReturnValues) > 0 {
		for _, v := range(command.ReturnValues) {
			vals = vals+v.(string)+" "
		}
		msg = vals
	}
	msg = "["+command.Name+" ->] "+status+" "+msg+"from "+command.From
	return msg
}

func PrintCommandReport(command *datatypes.Command) {
	fmt.Println(getTime(), GetCommandReportMsg(command))
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
