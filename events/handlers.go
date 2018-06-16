package events

import (
	"github.com/synw/microb/msgs"
	"github.com/synw/microb/services/logs"
	"github.com/synw/microb/types"
)

/*
Handles the print and log actions on the event
*/
func handle(event *types.Event) {
	// check if the event has to be logged
	log := true
	if event.Cmd != nil {
		if event.Cmd.NoLog == true {
			log = false
		}
		if event.Cmd.Trace != nil {
			event.Msg = event.Cmd.ErrMsg
		}
	}
	msgs.Event(event)
	if log == true {
		// log the event
		logs.Event(event)
	}
}
