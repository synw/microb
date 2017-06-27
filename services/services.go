package services

import (
	"errors"
	"fmt"
	color "github.com/acmacalister/skittles"
	"github.com/synw/microb/libmicrob/events"
	"github.com/synw/microb/libmicrob/types"
	"github.com/synw/terr"
)

var ValidCommands []string
var isDev = false

func GetService(name string) *types.Service {
	s := All
	if isDev == true {
		s = AllDev
	}
	serv := s[name]
	return serv
}

func InitServices(dev bool, verbosity int, servs []string) (map[string]*types.Service, []*terr.Trace) {
	if verbosity > 0 {
		fmt.Println("Initializing services ...")
	}
	services := make(map[string]*types.Service)
	var trs []*terr.Trace
	for _, name := range servs {
		if verbosity > 2 {
			fmt.Println("Initializing service", name)
		}
		service := GetService(name)
		tr := service.Init(dev, verbosity)
		if tr != nil {
			tr = terr.Add("initServices", errors.New("Unable to initialize "+name+" service"))
			msg := tr.Formatc()
			events.Err(name, "services.InitServices", msg, tr.ToErr())
			trs = append(trs, tr)
		} else {
			msg := color.BoldWhite(name) + " service initialized"
			events.State(name, "services.Init", msg, nil)
		}
		services[name] = service
		ValidCommands = append(ValidCommands, service.Cmds...)
	}
	return services, trs
}

func Dispatch(cmd *types.Command, c chan *types.Command) {
	service := GetService(cmd.Service)
	com := service.Dispatch(cmd)
	c <- com
}

// used by cli
func getAllValidCommands() []string {
	vc := []string{}
	for _, s := range All {
		vc = append(vc, All[s.Name].Cmds...)
	}
	return vc
}

func CmdIsValid(cmd *types.Command) bool {
	name := cmd.Name
	for _, c := range getAllValidCommands() {
		if c == name {
			return true
		}
	}
	return false
}
