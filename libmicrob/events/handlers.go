package events

import (
	"github.com/synw/microb/libmicrob/msgs"
	"github.com/synw/microb/libmicrob/types"
	"github.com/synw/microb/services/logs"
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
	}
	if log == true {
		// print event
		msgs.Event(event)
		// log event
		logs.Event(event)
	}
}
