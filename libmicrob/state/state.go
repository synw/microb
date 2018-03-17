package state

import (
	"fmt"
	color "github.com/acmacalister/skittles"
	"github.com/looplab/fsm"
	"github.com/synw/centcom"
	"github.com/synw/microb/libmicrob/conf"
	"github.com/synw/microb/libmicrob/msgs"
	"github.com/synw/microb/libmicrob/types"
	"github.com/synw/microb/services"
	"github.com/synw/terr"
)

type VerbState struct {
	State *fsm.FSM
}

var Verb = VerbState{}
var Server = &types.WsServer{}
var Cli *centcom.Cli
var Services map[string]*types.Service
var ValidCmds map[string]*types.Cmd
var Conf *types.Conf

func Init() *terr.Trace {
	fmt.Println("Starting Microb instance ...")
	// get config
	Conf, tr := conf.GetConf()
	if tr != nil {
		tr = terr.Pass("Init", tr)
		return tr
	}
	// get server
	Server, tr = conf.GetServer(Conf)
	if tr != nil {
		tr = terr.Pass("Init", tr)
		return tr
	}
	Cli, tr = initWsCli()
	if tr != nil {
		tr = terr.Pass("Init", tr)
		return tr
	}
	// get services
	Services, tr = getServices(Conf.Services)
	if tr != nil {
		tr = terr.Pass("Init", tr)
		return tr
	}
	msgs.Ready("Services are ready")
	return nil
}

func getServices(servs []string) (map[string]*types.Service, *terr.Trace) {
	srvs := make(map[string]*types.Service)
	manSrvs := services.Services
	for _, name := range servs {
		for k, v := range manSrvs {
			if k == name {
				srvs[k] = v
				msgs.Status("Initializing " + color.BoldWhite(v.Name) + " service")
				break
			}
		}
	}
	return srvs, nil
}

func initWsCli() (*centcom.Cli, *terr.Trace) {
	cli := centcom.NewClient(Server.Addr, Server.Key)
	err := centcom.Connect(cli)
	if err != nil {
		trace := terr.New("initWsCli", err)
		var cli *centcom.Cli
		return cli, trace
	}
	cli.IsConnected = true
	msgs.Ok("Websockets client connected")
	err = cli.CheckHttp()
	if err != nil {
		trace := terr.New("InitCli", err)
		return cli, trace
	}
	msgs.Ok("Websockets http transport ready")
	return cli, nil
}
