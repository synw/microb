package events

import (
	"fmt"
	color "github.com/acmacalister/skittles"
	"github.com/synw/microb/libmicrob/log"
	"github.com/synw/microb/libmicrob/state"
	"github.com/synw/microb/libmicrob/types"
)

func handle(event *types.Event) {
	if state.Verbosity > 0 {
		if event.Msg != "" {
			fmt.Println(getMsg(event))
		}
	}
	if event.Err != nil {
		if state.Verbosity > 0 {
			fmt.Println(event.Err.Error())
		}
	}
	log.New(event.Service, state.Conf.Name, event.Class, event.Msg)
}

func getMsg(event *types.Event) string {
	out := getFormatedMsgNoTime(event)
	return out
}

func getFormatedMsgNoTime(event *types.Event) string {
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
