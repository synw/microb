package libmicrob

import (
	"fmt"
	"github.com/looplab/fsm"
	"github.com/synw/microb/libmicrob/conf"
	"github.com/synw/microb/libmicrob/types"
)

type DevState struct {
	State *fsm.FSM
}

type VerbState struct {
	State *fsm.FSM
}

var Verb = VerbState{}
var Dev = DevState{}

func Verbose() bool {
	if Verb.State.Current() != "zero" {
		return true
	}
	return false
}

func Init(verb int, dev bool) *types.Conf {
	if verb > 0 {
		fmt.Println("Starting Microb instance ...")
	}
	initDev()
	initVerb()
	if dev == true {
		Dev.State.Event("activate")
	}
	if verb == 0 {
		Verb.State.Event("setZero")
	} else if verb == 1 {
		Verb.State.Event("setOne")
	} else if verb == 2 {
		Verb.State.Event("setTwo")
	}
	conf, tr := conf.GetConf(Dev.State.Current())
	if tr != nil {
		tr.Error()
	}
	return conf
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
