package services

import (
	//"fmt"
	//m "github.com/synw/microb/libmicrob"
	"github.com/synw/microb/libmicrob/types"
	"github.com/synw/terr"
)

func GetService(name string) *types.Service {
	return services[name]
}

func Init(servs []string, start bool) (map[string]*types.Service, *terr.Trace) {
	srv := make(map[string]*types.Service)
	return srv, nil
}
