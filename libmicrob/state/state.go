package state

import (
	color "github.com/acmacalister/skittles"
	"github.com/synw/centcom"
	config "github.com/synw/microb/libmicrob/conf"
	"github.com/synw/microb/libmicrob/msgs"
	"github.com/synw/microb/libmicrob/redis"
	"github.com/synw/microb/libmicrob/types"
	"github.com/synw/microb/services"
	"github.com/synw/microb/services/logs"
	"github.com/synw/terr"
)

func Init(dev bool, start bool) (*types.State, *terr.Trace) {
	state := &types.State{}
	msgs.Msg("Starting Microb instance ...")
	// get config
	//var tr *terr.Trace
	conf, tr := config.GetConf()
	state.Conf = conf
	if tr != nil {
		tr = terr.Pass("Init", tr)
		return state, tr
	}
	// get server
	state.WsServer, tr = config.GetServer(state.Conf)
	if tr != nil {
		tr = terr.Pass("Init", tr)
		return state, tr
	}
	state.Cli, tr = initWsCli(state)
	if tr != nil {
		tr = terr.Pass("Init", tr)
		return state, tr
	}
	// initialize Redis connection
	tr = redis.InitRedis(conf)
	if tr != nil {
		tr.Fatal()
	}
	// initialize logger
	logs.Init(conf)
	// get services
	state.Services, tr = initServices(state.Conf.Services, dev, start)
	if tr != nil {
		tr = terr.Pass("Init", tr)
		return state, tr
	}
	msgs.Ready("Services are ready")
	// get commands
	state.Cmds = make(map[string]*types.Cmd)
	for _, srv := range state.Services {
		for cname, cmd := range srv.Cmds {
			state.Cmds[srv.Name+"_"+cname] = cmd
		}
	}
	return state, nil
}

func initServices(servs []string, dev bool, start bool) (map[string]*types.Service, *terr.Trace) {
	srvs := make(map[string]*types.Service)
	manSrvs := services.Services
	for _, name := range servs {
		for k, srv := range manSrvs {
			if k == name {
				srvs[k] = srv
				msgs.Status("Initializing " + color.BoldWhite(srv.Name) + " service")
				err := srv.Init(dev, start)
				if err != nil {
					tr := terr.New("state.getServices", err)
					return srvs, tr
				}
				break
			}
		}
	}
	return srvs, nil
}

func initWsCli(state *types.State) (*centcom.Cli, *terr.Trace) {
	cli := centcom.NewClient(state.WsServer.Addr, state.WsServer.Key)
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
