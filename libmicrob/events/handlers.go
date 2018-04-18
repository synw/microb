package events

import (
	"github.com/synw/microb/libmicrob/msgs"
	"github.com/synw/microb/libmicrob/types"
	"github.com/synw/microb/services/logs"
)

func handle(event *types.Event) {
	log := true
	if event.Cmd != nil {
		if event.Cmd.NoLog == true {
			log = false
		}
	}
	if log == true {
		logs.Event(event)
		msgs.PrintEvent(event)
	}
}
