package services

import (
	http "github.com/synw/microb-http/manifest"
	"github.com/synw/microb/libmicrob/types"
)

var All = map[string]*types.Service{
	"http": http.Service,
}
