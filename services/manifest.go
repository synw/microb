package services

import (
	//http "github.com/synw/microb-http/manifest"
	infos "github.com/synw/microb-infos/manifest"
	"github.com/synw/microb/libmicrob/types"
)

var services = map[string]*types.Service{
	"infos": infos.Service,
	//"http": http.Service,
}
