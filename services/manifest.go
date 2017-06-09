package services

import (
	grgCmd "github.com/synw/microb-goregraph/cmd"
	grgState "github.com/synw/microb-goregraph/state"
	infoCmd "github.com/synw/microb/libmicrob/cmd/info"
	"github.com/synw/microb/libmicrob/datatypes"
	httpCmd "github.com/synw/microb/services/httpServer/cmd"
	httpState "github.com/synw/microb/services/httpServer/state"
)

var infoCmds = []string{"ping"}
var httpCmds = []string{"start", "stop"}

var grgCmds = []string{"start", "stop"}

var All = map[string]*datatypes.Service{
	"info":      New("info", infoCmds),
	"http":      New("http", httpCmds),
	"goregraph": New("goregraph", grgCmds),
}

var initState = map[string]interface{}{
	"http":      httpState.InitState,
	"goregraph": grgState.InitState,
}
var initDispatch = map[string]interface{}{
	"http":      httpCmd.Dispatch,
	"info":      infoCmd.Dispatch,
	"goregraph": grgCmd.Dispatch,
}
