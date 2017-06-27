package services

import (
	http "github.com/synw/microb-http/manifest"
	"github.com/synw/microb/libmicrob/types"
	"github.com/synw/microb/services/info"
)

var All = map[string]*types.Service{
	"info": info.Service,
	"http": http.Service,
}
