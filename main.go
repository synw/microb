package main

import (
	"errors"
	"flag"
	"fmt"
	color "github.com/acmacalister/skittles"
	"github.com/synw/centcom"
	"github.com/synw/microb/libmicrob/cmd"
	"github.com/synw/microb/libmicrob/events"
	"github.com/synw/microb/libmicrob/log"
	"github.com/synw/microb/libmicrob/state"
	"github.com/synw/microb/services"
	"github.com/synw/terr"
)

var dev = flag.Bool("d", false, "Dev mode")
var verbosity = flag.Int("v", 1, "Verbosity")
var start = flag.Bool("s", false, "Start services")

func main() {
	if state.Verbosity > 0 {
		fmt.Println("Starting Microb service ...")
	}
	flag.Parse()
	// init state
	conf, tr := state.InitState(*dev, *verbosity)
	if tr != nil {
		err := errors.New("Unable to initialize state")
		tr = terr.Add("main", err, tr)
		events.Err("microb", "main", tr.Formatc(), tr.ToErr())
		return
	}
	// init logger
	log.Init(state.Cli)
	if *dev == true {
		events.State("main", "state.InitState", "Dev mode is on", nil)
	}
	events.State("main", "state.InitState", "Commands transport layer operational", nil)
	defer centcom.Disconnect(state.Cli)
	if state.Verbosity > 2 {
		terr.Ok("Initialized state")
	}
	// init services
	services, trs := services.InitServices(*dev, *verbosity, conf.Services, *start)
	state.Services = services
	if trs != nil {
		for _, tr = range trs {
			tr.Printf("unable to inialize service")
		}
	}
	events.Ready("microb", "main", "Services are ready", nil)
	// connect on the commands channel
	err := state.Cli.Subscribe(state.Server.CmdChanIn)
	if err != nil {
		fmt.Println(err)
	}
	// listen
	go func() {
		msg := color.BoldWhite("Ready") + ": listening for commands at " + state.Cli.Addr + ":" + " on channel " + state.Server.CmdChanIn + " ..."
		events.State("microb", "main", msg, nil)
		for msg := range state.Cli.Channels {
			if msg.Channel == state.Server.CmdChanIn {
				//fmt.Println("PAYLOAD", msg.Payload.(map[string]interface{}))
				cmd.Run(msg.Payload, state.Cli, state.Server)
			}
		}
	}()
	// idle
	select {}
}
