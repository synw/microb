package mutate

import (
	"errors"
	"github.com/synw/microb/services/httpServer"
	"github.com/synw/microb/services/httpServer/state"
	"github.com/synw/terr"
)

func StartHttpServer() *terr.Trace {
	if state.HttpServer.Running == true {
		err := errors.New("Http server is already running")
		tr := terr.New("state.mutate.StartHttpServer", err)
		return tr
	}
	go httpServer.InitHttpServer(true)
	return nil
}

func StopHttpServer() *terr.Trace {
	if state.HttpServer.Running == false {
		err := errors.New("Http server is not running")
		tr := terr.New("state.mutate.StopHttpServer", err)
		return tr
	}
	tr := httpServer.Stop()
	if tr != nil {
		return tr
	}
	return nil
}
