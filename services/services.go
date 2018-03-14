package services

import (
	color "github.com/acmacalister/skittles"
	m "github.com/synw/microb/libmicrob"
	"github.com/synw/microb/libmicrob/types"
	"github.com/synw/terr"
)

func GetService(name string) *types.Service {
	return services[name]
}

func Init(servs []string, start bool) (map[string]*types.Service, *terr.Trace) {
	srv := make(map[string]*types.Service)
	srv["infos"] = GetService("infos")
	for _, name := range servs {
		s := GetService(name)
		s.Init(m.Verbose())
		m.Status("Initializing " + color.BoldWhite(s.Name) + " service")
	}
	return srv, nil
}

func Dispatch(cmd *types.Cmd, c chan *types.Cmd) {
	com := cmd.Exec.(func(*types.Cmd) *types.Cmd)(cmd)
	c <- com
}
