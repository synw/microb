package main

import (
	"fmt"
    "flag"
    "errors"
    "strconv"
    "github.com/synw/centcom"
    "github.com/synw/terr"
    "github.com/synw/microb/libmicrob/state"
    "github.com/synw/microb/libmicrob/cmd"
    "github.com/synw/microb/libmicrob/events"
    "github.com/synw/microb/libmicrob/httpServer"
    "github.com/synw/microb/libmicrob/db"
)


var dev_mode = flag.Bool("d", false, "Dev mode")
var verbosity = flag.Int("v", 1, "Verbosity")

func main() {
	flag.Parse()
	tr := state.InitState(*dev_mode, *verbosity)
	if tr != nil {
		err := errors.New("Unable to initialize state")
		tr = terr.Add("main", err, tr)
		events.Error(tr)
		return
	}
	defer centcom.Disconnect(state.Cli)
	if state.Verbosity > 2 {
		fmt.Println(terr.Ok("Initialized state"))
	}
	// init http server
	httpServer.InitHttpServer()
	// init database
	name := "normal"
	if *dev_mode == true {
		name = "dev"
	}
	tr = db.InitDb(name)
	if tr != nil {
		tr.Formatc()
	}
	// connect on the commands channel
	err := state.Cli.Subscribe(state.Server.CmdChanIn)
	if err != nil {
		fmt.Println(err)
	}
	// listen
	go func() {
		m := "Listening for commands at "+state.Cli.Host+":"+strconv.Itoa(state.Cli.Port)+" on channel "+state.Server.CmdChanIn+" ..."
		fmt.Println(m)
		for msg := range(state.Cli.Channels) {
			if msg.Channel == state.Server.CmdChanIn {
				//fmt.Println("PAYLOAD", msg.Payload.(map[string]interface{}))
				cmd.Run(msg.Payload)
			}
		}
	}()
	// idle
	select {}
}
