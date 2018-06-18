package events

import (
	"github.com/synw/microb/msgs"
	"github.com/synw/microb/services/logs"
	"github.com/synw/microb/types"
)

// Handles the print and log actions on the event
func handle(event *types.Event) {
	// check if the event has to be logged
	log := true
	if event.Cmd != nil {
		if event.Cmd.NoLog == true {
			log = false
		}
	}
	// print the event
	msgs.Msg("EVENT MSG")
	msgs.Event(event)
	msgs.Msg("END EVENT MSG")
	// log the event
	if log == true {
		logs.Event(event)
	}
}
