package libmicrob

import (
	"github.com/looplab/fsm"
)

func (verb *VerbState) Mutate(e *fsm.Event) {
	State("Verbosity set to " + e.Dst)
}
