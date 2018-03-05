package libmicrob

import (
	"fmt"
	"github.com/looplab/fsm"
	"github.com/synw/centcom"
	"github.com/synw/microb/libmicrob/conf"
	"github.com/synw/microb/libmicrob/types"
	"github.com/synw/terr"
)

type VerbState struct {
	State *fsm.FSM
}

var Verb = VerbState{}
var server = &types.WsServer{}
var cli *centcom.Cli

func Verbose() bool {
	if Verb.State.Current() == "one" {
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
	} else {
		Verb.State.Event("setOne")
	}
	config, tr := conf.GetConf()
	if tr != nil {
		tr = terr.Pass("Init", tr)
		return config, tr
	}
	server, tr = conf.GetServer(config)
	if tr != nil {
		tr = terr.Pass("Init", tr)
		return config, tr
	}
	cli, tr = initWsCli()
	if tr != nil {
		tr = terr.Pass("Init", tr)
		return config, tr
	}
	return config, nil
}

func initVerb() {
	Verb.State = fsm.NewFSM(
		"one",
		fsm.Events{
			{Name: "setZero", Src: []string{"zero", "one"}, Dst: "zero"},
			{Name: "setOne", Src: []string{"zero", "one"}, Dst: "one"},
		},
		fsm.Callbacks{
			"enter_state": func(e *fsm.Event) { Verb.Mutate(e) },
		},
	)
}

func initWsCli() (*centcom.Cli, *terr.Trace) {
	cli := centcom.NewClient(server.Addr, server.Key)
	err := centcom.Connect(cli)
	if err != nil {
		trace := terr.New("initWsCli", err)
		var cli *centcom.Cli
		return cli, trace
	}
	cli.IsConnected = true
	Ok("Websockets client connected")
	err = cli.CheckHttp()
	if err != nil {
		trace := terr.New("InitCli", err)
		return cli, trace
	}
	Ok("Websockets http transport ready")
	return cli, nil
}
