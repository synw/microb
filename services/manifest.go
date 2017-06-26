package services

import (
	dashboardCmd "github.com/synw/microb-dashboard/cmd"
	dashboardState "github.com/synw/microb-dashboard/state"
	httpCmd "github.com/synw/microb-http/cmd"
	httpState "github.com/synw/microb-http/state"
	infoCmd "github.com/synw/microb/libmicrob/cmd/info"
	"github.com/synw/microb/libmicrob/types"
)

var infoCmds = []string{"ping"}
var httpCmds = []string{"start", "stop"}
var dashboardCmds = []string{"start", "stop"}

var All = map[string]*types.Service{
	"info":      New("info", infoCmds),
	"http":      New("http", httpCmds),
	"dashboard": New("dashboard", dashboardCmds),
}

var initState = map[string]interface{}{
	"http":      httpState.InitState,
	"dashboard": dashboardState.InitState,
}
var dispatch = map[string]interface{}{
	"http":      httpCmd.Dispatch,
	"info":      infoCmd.Dispatch,
	"dashboard": dashboardCmd.Dispatch,
}
