package services

import (
	cmdInfo "github.com/synw/microb/libmicrob/cmd/info"
	cmdHttp "github.com/synw/microb/services/httpServer/cmd"
	httpState "github.com/synw/microb/services/httpServer/state"
)

var initState = map[string]interface{}{
	"http": httpState.InitState,
}
var initDispatch = map[string]interface{}{
	"http": cmdHttp.Dispatch,
	"info": cmdInfo.Dispatch,
}
