package infos

import (
	"github.com/synw/microb/types"
)

var Service *types.Service = &types.Service{
	"infos",
	getCmds(),
	initService,
}
