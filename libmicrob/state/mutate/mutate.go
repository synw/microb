package mutate

import (
	"errors"
	"github.com/synw/terr"
	"github.com/synw/microb/libmicrob/httpServer"
	"github.com/synw/microb/libmicrob/state"
	"github.com/synw/microb/libmicrob/events"
)


func StartHttpServer() *terr.Trace {
	if state.HttpServer.Running == true {
		err := errors.New("Http server is already running")
		tr := terr.New("state.mutate.StartHttpServer", err)
		return tr
	}
	go httpServer.Run()
	return nil
}

func StopHttpServer() {
	tr := httpServer.Stop()
	if tr != nil {
		events.Err("error", "state.mutate.StopHttpServer", tr)
	}
	return
}
