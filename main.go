package main

import (
	"fmt"
    "flag"
    "errors"
    "strconv"
    "github.com/synw/centcom"
    "github.com/synw/terr"
    "github.com/synw/microb/libmicrob/state"
    "github.com/synw/microb/libmicrob/events"
)


var dev_mode = flag.Bool("d", false, "Dev mode")
var verbosity = flag.Int("v", 1, "Verbosity")

func main() {
	flag.Parse()
	cli, trace := state.InitState(*dev_mode, *verbosity)
	if trace != nil {
		err := errors.New("Unable to initialize state")
		trace = terr.Add("main", err, trace)
		events.Error(trace)
		return
	}
	defer centcom.Disconnect(cli)
	if state.Verbosity > 2 {
		fmt.Println(terr.Ok("Initialized state"))
	}
	// connect on the commands channel
	err := cli.Subscribe(state.Server.CmdChannel)
	if err != nil {
		fmt.Println(err)
	}
	// listen
	go func() {
		msg := "Listening for commands at "+cli.Host+":"+strconv.Itoa(cli.Port)+" on channel "+state.Server.CmdChannel+"..."
		fmt.Println(msg)
		for msg := range(cli.Channels) {
			if msg.Channel == state.Server.CmdChannel {
				fmt.Println("PAYLOAD", msg.Payload, msg.UID)
			}
		}
	}()
	// idle
	for {
		
	}
}
