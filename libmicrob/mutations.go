package libmicrob

import (
	"github.com/looplab/fsm"
)

func (dev *DevState) Mutate(e *fsm.Event) {
	if e.Event == "activate" {
		state("Dev mode activated")
	} else {
		state("Dev mode deactivated")
	}
}

func (verb *VerbState) Mutate(e *fsm.Event) {
	state("Verbosity set to " + e.Dst)
}
