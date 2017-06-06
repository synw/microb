package services

import (
	cmdInfo "github.com/synw/microb/libmicrob/cmd/info"
	"github.com/synw/microb/libmicrob/datatypes"
	cmdHttp "github.com/synw/microb/services/httpServer/cmd"
	httpState "github.com/synw/microb/services/httpServer/state"
)

var infoCmds = []string{"ping"}
var httpCmds = []string{"start", "stop"}

var All = map[string]*datatypes.Service{
	"info": New("info", infoCmds),
	"http": New("http", httpCmds),
}

var initState = map[string]interface{}{
	"http": httpState.InitState,
}
var initDispatch = map[string]interface{}{
	"http": cmdHttp.Dispatch,
	"info": cmdInfo.Dispatch,
}
