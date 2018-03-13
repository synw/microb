package main

import (
	"flag"
	"fmt"
	color "github.com/acmacalister/skittles"
	m "github.com/synw/microb/libmicrob"
	"github.com/synw/microb/libmicrob/cmd"
	events "github.com/synw/microb/libmicrob/events"
	"github.com/synw/microb/services"
)

var dev = flag.Bool("d", false, "Dev mode")
var verb = flag.Int("v", 1, "Verbosity")

func main() {
	flag.Parse()
	conf, tr := m.Init(*verb, *dev)
	if tr != nil {
		tr.Print()
		return
	}
	m.Ok("State initialized")
	// TODO : start services flag
	services.Init(conf.Services, true)
	m.Ready("Services are ready")
	// connect on the commands channel
	err := m.Cli.Subscribe(m.Server.CmdChanIn)
	if err != nil {
		fmt.Println(err)
	}
	// listen
	go func() {
		msg := color.BoldWhite("Ready") + ": listening for commands at " + m.Cli.Addr + ":" + " on channel " + m.Server.CmdChanIn + " ..."
		events.State(msg)
		for msg := range m.Cli.Channels {
			if msg.Channel == m.Server.CmdChanIn {
				//m.Debug(msg.Payload.(map[string]interface{}))
				cmd.Run(msg.Payload)
			}
		}
	}()
	// idle
	select {}
}
