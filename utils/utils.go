package utils

import (
	"fmt"
	"time"
	"github.com/acmacalister/skittles"
)

var info_m = " ["+skittles.Green("info")+"]"
var event_m = " ["+skittles.Yellow("event")+"]"
var command_m = " [=> "+skittles.Cyan("command")+"]"
var error_m = " ["+skittles.BoldRed("error")+"]"
var metric_m = " ["+skittles.Cyan("metric")+"]"

func GetTime() string {
	t := time.Now()
	return t.Format("15:04:05")
}

func GetEvent(itype string, msg string) string {
	event := GetTime()
	if (itype == "info") {
		event = event+info_m+" "+msg
	} else if (itype == "event") {
		event = event+event_m+" "+msg
	} else if (itype == "command") {
		event = event+command_m+" "+msg
	} else if (itype == "error") {
		event = event+error_m+" "+msg
	} else if (itype == "metric") {
		event = event+metric_m+" "+msg
	} else if (itype == "nil") {
		event = event+" "+msg
	}
	return event
}

func PrintEvent(itype string, msg string) {
	go fmt.Println(GetEvent(itype, msg))
}