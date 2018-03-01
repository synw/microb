package state

import (
	"github.com/looplab/fsm"
	"github.com/synw/microb/libmicrob/events"
)

func (dev *DevState) Mutate(e *fsm.Event) {
	if e.Event == "activate" {
		args := make(map[string]interface{})
		args["class"] = "state"
		events.New("Dev mode activated", args)
	} else {
		events.New("Dev mode deactivated")
	}
}
