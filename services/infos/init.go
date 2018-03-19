package infos

import (
	"github.com/synw/microb/libmicrob/types"
)

var Service *types.Service = &types.Service{
	"infos",
	getCmds(),
	initService,
}
