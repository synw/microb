package state

import (
	"github.com/looplab/fsm"
)

var Verb = VerbState{}
var Dev = DevState{}

func Init(verb int, isDev bool) {
	initDev()
	if isDev == true {
		Dev.State.Event("activate")
	}
}

func initDev() {
	Dev.State = fsm.NewFSM(
		"off",
		fsm.Events{
			{Name: "activate", Src: []string{"off"}, Dst: "on"},
			{Name: "deactivate", Src: []string{"on"}, Dst: "off"},
		},
		fsm.Callbacks{
			"enter_state": func(e *fsm.Event) { Dev.Mutate(e) },
		},
	)
}
