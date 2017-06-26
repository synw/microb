package services

import (
	//grgCmd "github.com/synw/microb-goregraph/cmd"
	//grgState "github.com/synw/microb-goregraph/state"
	dashboardCmd "github.com/synw/microb-dashboard/cmd"
	dashboardState "github.com/synw/microb-dashboard/state"
	infoCmd "github.com/synw/microb/libmicrob/cmd/info"
	"github.com/synw/microb/libmicrob/types"
	httpCmd "github.com/synw/microb/services/httpServer/cmd"
	httpState "github.com/synw/microb/services/httpServer/state"
)

var infoCmds = []string{"ping"}
var httpCmds = []string{"start", "stop"}
var dashboardCmds = []string{"start", "stop"}

//var grgCmds = []string{"start", "stop"}

var All = map[string]*types.Service{
	"info":      New("info", infoCmds),
	"http":      New("http", httpCmds),
	"dashboard": New("dashboard", dashboardCmds),
	//"goregraph": New("goregraph", grgCmds),
}

var initState = map[string]interface{}{
	"http":      httpState.InitState,
	"dashboard": dashboardState.InitState,
	//"goregraph": grgState.InitState,
}
var initDispatch = map[string]interface{}{
	"http":      httpCmd.Dispatch,
	"info":      infoCmd.Dispatch,
	"dashboard": dashboardCmd.Dispatch,
	//"goregraph": grgCmd.Dispatch,
}
