package main

import (
	"flag"
	//"github.com/synw/microb/libmicrob/msgs"
	"github.com/synw/microb/libmicrob/state"
)

var dev = flag.Bool("d", false, "Dev mode")
var verbosity = flag.Int("v", 1, "Verbosity")

func main() {
	flag.Parse()
	state.Init(*verbosity, *dev)
	/*if state.Verb > 0 {
		msgs.Print("Microb state initialized")
		if state.Dev == true {
			msgs.Print("Dev mode is on")
		}
	}*/
}
