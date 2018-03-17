package main

import (
	"flag"
	"fmt"
	"github.com/synw/microb/libmicrob/cmds"
	"github.com/synw/microb/libmicrob/events"
	"github.com/synw/microb/libmicrob/msgs"
	"github.com/synw/microb/libmicrob/state"
)

func main() {
	flag.Parse()
	tr := state.Init()
	if tr != nil {
		tr.Print()
		return
	}
	msgs.Ok("State initialized")
	// connect on the commands channel
	err := state.Cli.Subscribe(state.Server.CmdChanIn)
	if err != nil {
		fmt.Println(err)
	}
	// listen
	go func() {
		msg := "listening for commands at " + state.Cli.Addr + ":" + " on channel " + state.Server.CmdChanIn + " ..."
		events.State(msg)
		for msg := range state.Cli.Channels {
			if msg.Channel == state.Server.CmdChanIn {
				//m.Debug(msg.Payload.(map[string]interface{}))
				cmds.Run(msg.Payload)
			}
		}
	}()
	// idle
	select {}
}
