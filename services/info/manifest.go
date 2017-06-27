package info

import (
	"github.com/synw/microb/libmicrob/cmd/info"
	"github.com/synw/microb/libmicrob/types"
	"github.com/synw/terr"
)

var Service *types.Service = &types.Service{
	"info",
	[]string{"ping"},
	ini,
	dispatch,
}

func ini(dev bool, verbosity int) *terr.Trace {
	return nil
}

func dispatch(cmd *types.Command) *types.Command {
	return info.Dispatch(cmd)
}
