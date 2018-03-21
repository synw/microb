package main

import (
	"flag"
	"fmt"
	"github.com/synw/microb/libmicrob/cmds"
	"github.com/synw/microb/libmicrob/events"
	"github.com/synw/microb/libmicrob/msgs"
	"github.com/synw/microb/libmicrob/state"
)

var dev = flag.Bool("d", false, "Dev mode")
var start = flag.Bool("s", false, "Start the services")

func main() {
	flag.Parse()
	state, tr := state.Init(*dev, *start)
	if tr != nil {
		tr.Print()
		return
	}
	msgs.Ok("State initialized")
	// connect on the commands channel
	err := state.Cli.Subscribe(state.WsServer.CmdChanIn)
	if err != nil {
		fmt.Println(err)
	}
	// listen
	go func() {
		msg := "listening for commands at " + state.Cli.Addr + ":" + " on channel " + state.WsServer.CmdChanIn + " ..."
		events.State(msg)
		for msg := range state.Cli.Channels {
			if msg.Channel == state.WsServer.CmdChanIn {
				//msgs.Debug(msg.Payload.(map[string]interface{}))
				cmds.Run(msg.Payload, state)
			}
		}
	}()
	// idle
	select {}
}
