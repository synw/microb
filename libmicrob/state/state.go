package state

import (
	"fmt"
	"github.com/synw/centcom"
	"github.com/synw/microb/libmicrob/conf"
	"github.com/synw/microb/libmicrob/datatypes"
	"github.com/synw/microb/services"
	httpConf "github.com/synw/microb/services/httpServer/conf"
	"github.com/synw/terr"
)

var Debug = true
var Server = &datatypes.Server{}
var Verbosity int = 1
var Cli *centcom.Cli
var DocDb = &datatypes.Database{}
var Conf map[string]interface{}
var Dev bool
var Services []*datatypes.Service
var ValidCommands = getValidCommands()

func getValidCommands() []string {
	httpCmds := httpConf.ValidCommands
	infoCmds := []string{"ping"}
	validCommands := []string{}
	allComs := append(httpCmds, infoCmds...)
	validCommands = append(validCommands, allComs...)
	return validCommands
}

func InitState(dev bool, verbosity int) *terr.Trace {
	Verbosity = verbosity
	Dev = dev
	// command channel server
	server, trace := conf.GetServer(dev)
	if trace != nil {
		trace = terr.Pass("stateInit.State", trace)
		return trace
	}
	Server = server
	// Websockets client
	if Verbosity > 0 {
		fmt.Println("Initializing commands transport layer ...")
	}
	cli, tr := initWsCli()
	if tr != nil {
		tr = terr.Pass("state.InitState", tr)
		return tr
	}
	Cli = cli
	// conf
	Conf, tr = conf.GetConf(dev)
	if tr != nil {
		return tr
	}
	// services
	initServices()
	return nil
}

// internal methods

func initServices() {
	if Verbosity > 0 {
		fmt.Println("Initializing services ...")
	}
	for _, el := range Conf["services"].([]interface{}) {
		name := el.(string)
		s := services.New(name)
		Services = append(Services, s)
		if Verbosity > 1 {
			fmt.Println("Registering service", name)
		}
	}
}

func initWsCli() (*centcom.Cli, *terr.Trace) {
	cli := centcom.NewClient(Server.WsHost, Server.WsPort, Server.WsKey)
	err := centcom.Connect(cli)
	if err != nil {
		trace := terr.New("InitCli", err)
		var cli *centcom.Cli
		return cli, trace
	}
	cli.IsConnected = true
	if Verbosity > 1 {
		fmt.Println(terr.Ok("Websockets client connected"))
	}
	err = cli.CheckHttp()
	if err != nil {
		trace := terr.New("InitCli", err)
		return cli, trace
	}
	if Verbosity > 1 {
		fmt.Println(terr.Ok("Websockets http transport ready"))
	}
	return cli, nil
}
