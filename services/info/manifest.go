package info

import (
	"github.com/synw/microb/libmicrob/types"
	"github.com/synw/terr"
)

var Service *types.Service = &types.Service{
	"info",
	[]string{"ping"},
	ini,
	dispatch,
}

func ini(dev bool, verb int, start bool) *terr.Trace {
	return nil
}

func dispatch(cmd *types.Cmd) *types.Cmd {
	return Dispatch(cmd)
}
