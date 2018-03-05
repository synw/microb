package libmicrob

import (
	"fmt"
	"github.com/looplab/fsm"
	"github.com/synw/microb/libmicrob/conf"
	"github.com/synw/microb/libmicrob/types"
	"github.com/synw/terr"
)

type VerbState struct {
	State *fsm.FSM
}

var Verb = VerbState{}

func Verbose() bool {
	if Verb.State.Current() != "zero" {
		return true
	}
	return false
}

func Init(verb int, dev bool) (*types.Conf, *terr.Trace) {
	if verb > 0 {
		fmt.Println("Starting Microb instance ...")
	}
	initVerb()
	if verb == 0 {
		Verb.State.Event("setZero")
	} else if verb == 1 {
		Verb.State.Event("setOne")
	} else if verb == 2 {
		Verb.State.Event("setTwo")
	}
	conf, tr := conf.GetConf()
	if tr != nil {
		terr.Pass("Init", tr)
		return conf, tr
	}
	return conf, nil
}

func initVerb() {
	Verb.State = fsm.NewFSM(
		"one",
		fsm.Events{
			{Name: "setZero", Src: []string{"zero", "one", "two"}, Dst: "zero"},
			{Name: "setOne", Src: []string{"zero", "one", "two"}, Dst: "one"},
			{Name: "setTwo", Src: []string{"zero", "one", "two"}, Dst: "two"},
		},
		fsm.Callbacks{
			"enter_state": func(e *fsm.Event) { Verb.Mutate(e) },
		},
	)
}
