package main

import (
	"fmt"
    "flag"
    "github.com/synw/terr"
    "github.com/synw/microb/libmicrob/state"
    "github.com/synw/microb/libmicrob/events"
)


var dev_mode = flag.Bool("d", false, "Dev mode")
var verbosity = flag.Int("v", 1, "Verbosity")

func main() {
	flag.Parse()
	trace := state.InitState(*dev_mode, *verbosity)
	if trace != nil {
		trace = terr.Pass("main", trace)
		events.Error(trace)
	} else {
		if state.Verbosity > 2 {
			terr.Ok("Initialized state")
		}
	}
	if state.Verbosity > 0 {
			fmt.Println("waiting...")
		}	
	for {
		
	}
}
