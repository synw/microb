package format

import (
	"time"
	"github.com/acmacalister/skittles"
	"github.com/synw/microb/libmicrob/datatypes"
	"github.com/synw/microb/libmicrob/events/format/methods"
)


func GetFormatedMsg(event_class string, event *datatypes.Event) string {
	t := GetTime()
	out := t+" "+getFormatedMsgNoTime(event_class, event)
	return out
}

func GetCommandReportMsg(command *datatypes.Command) string {
	report := getCommandReportMsg(command, true)
	return report
}

func GetFormatedCommandReportMsg(command *datatypes.Command) string {
	report := getCommandReportMsg(command, true)
	now := GetTime()
	msg := now+" "+report
	return msg
}

func GetFormatedCommandReportMsgSimple(command *datatypes.Command) string {
	report := getCommandReportMsg(command, false)
	now := GetTime()
	msg := now+" "+report
	return msg
}

func GetCommandReportMsgSimple(command *datatypes.Command) string {
	report := getCommandReportMsg(command, false)
	return report
}

func ErrorFormated(err error) string {
	var d map[string]interface{}
	from := ""
	event := &datatypes.Event{"error", from, err.Error(), d}
	fm := getFormatedMsgNoTime("error", event)
	return fm
}

func GetTime() string {
	t := time.Now()
	return t.Format("15:04:05")
}

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

func getFormatedMsgNoTime(event_class string, event *datatypes.Event) string {
	var out string
	msg := event.Message
	output_flags := getEventOutputFlags()
	out = out+output_flags[event_class]
	if (event_class == "request" || event_class == "request_error") {
		sc := event.Data["status_code"].(int)
		out = out+methods.FormatStatusCode(sc)
	}
	out = out+" "+msg
	return out
}

func getCommandReportMsg(command *datatypes.Command, from bool) string {
	var status string
	var vals string
	var msg string
	if command.Status == "success" {
		status = skittles.Green("ok")	
	} else if command.Status == "error" {
		status = skittles.Red("error")
	} else {
		status = command.Status
	}
	if len(command.ReturnValues) > 0 {
		for i, v := range(command.ReturnValues) {
			vals = vals+v.(string)
			if i < len(vals) {
				vals=vals+" "
			}
		}
		msg = vals
	}
	msg = "["+command.Name+" ->] "+status+" "+msg
	if from == true {
		msg = msg+"(from "+command.From+")"
	}
	return msg
}
