package main

import (
	"flag"
	m "github.com/synw/microb/libmicrob"
	"github.com/synw/microb/services"
)

var dev = flag.Bool("d", false, "Dev mode")
var verb = flag.Int("v", 1, "Verbosity")

func main() {
	flag.Parse()
	conf := m.Init(*verb, *dev)
	m.Ok("Microb state initialized")
	// TODO : start services flag
	services.Init(conf.Services, true)
}
