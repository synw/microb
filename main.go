package main

import (
	"fmt"
    "flag"
    "errors"
    "github.com/synw/terr"
    "github.com/synw/microb/libmicrob/state"
    "github.com/synw/microb/libmicrob/events"
)


var dev_mode = flag.Bool("d", false, "Dev mode")
var verbosity = flag.Int("v", 1, "Verbosity")

func main() {
	flag.Parse()
	_, trace := state.InitState(*dev_mode, *verbosity)
	if trace != nil {
		err := errors.New("Unable to initialize state")
		trace = terr.Add("main", err, trace)
		events.Error(trace)
		return
	}
	if state.Verbosity > 2 {
		fmt.Println(terr.Ok("Initialized state"))
	}
	if state.Verbosity > 0 {
		fmt.Println("waiting...")
	}	
	for {
		
	}
}
