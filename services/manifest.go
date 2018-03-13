package services

import (
	//http "github.com/synw/microb-http/manifest"
	"github.com/synw/microb/libmicrob/types"
	"github.com/synw/microb/services/infos"
)

var services = map[string]*types.Service{
	"infos": infos.Service,
	//"http": http.Service,
}
