package state

import (
	"github.com/looplab/fsm"
)

type DevState struct {
	State *fsm.FSM
}

type VerbState struct {
	State *fsm.FSM
}
