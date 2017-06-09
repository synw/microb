package main

import (
	"errors"
	"flag"
	"fmt"
	color "github.com/acmacalister/skittles"
	"github.com/synw/centcom"
	"github.com/synw/microb/libmicrob/cmd"
	"github.com/synw/microb/libmicrob/events"
	"github.com/synw/microb/libmicrob/state"
	"github.com/synw/microb/services"
	"github.com/synw/terr"
	"strconv"
)

var dev = flag.Bool("d", false, "Dev mode")
var verbosity = flag.Int("v", 1, "Verbosity")
var serve = flag.Bool("s", false, "Start http server")

func main() {
	if state.Verbosity > 0 {
		fmt.Println("Starting Microb service ...")
	}
	flag.Parse()
	// init state
	tr := state.InitState(*dev, *verbosity)
	if tr != nil {
		err := errors.New("Unable to initialize state")
		tr = terr.Add("main", err, tr)
		events.Err("microb", "main", tr.Formatc(), tr.ToErr())
		return
	}
	if *dev == true {
		events.State("main", "state.InitState", "Dev mode is on", nil)
	}
	events.State("main", "state.InitState", "Commands transport layer operational", nil)
	defer centcom.Disconnect(state.Cli)
	if state.Verbosity > 2 {
		fmt.Println(terr.Ok("Initialized state"))
	}
	// init services
	tr = services.InitServices(*dev, *verbosity)
	if tr != nil {
		tr.Fatal("Problem initilizing services")
	}
	events.Ready("microb", "main", "Services are ready", nil)
	// connect on the commands channel
	err := state.Cli.Subscribe(state.Server.CmdChanIn)
	if err != nil {
		fmt.Println(err)
	}
	// listen
	go func() {
		msg := color.BoldWhite("Ready") + ": listening for commands at " + state.Cli.Host + ":" + strconv.Itoa(state.Cli.Port) + " on channel " + state.Server.CmdChanIn + " ..."
		events.State("microb", "main", msg, nil)
		for msg := range state.Cli.Channels {
			if msg.Channel == state.Server.CmdChanIn {
				//fmt.Println("PAYLOAD", msg.Payload.(map[string]interface{}))
				cmd.Run(msg.Payload)
			}
		}
	}()
	// idle
	select {}
}
