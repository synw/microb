package state

import (
	"errors"
	"fmt"
	"github.com/synw/centcom"
	"github.com/synw/microb/libmicrob/conf"
	"github.com/synw/microb/libmicrob/datatypes"
	"github.com/synw/microb/services"
	httpConf "github.com/synw/microb/services/httpServer/conf"
	httpState "github.com/synw/microb/services/httpServer/state"
	"github.com/synw/terr"
)

var Debug = true
var Verbosity int = 1
var Cli *centcom.Cli
var Server = &datatypes.Server{}
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
	tr = initServices(dev, verbosity)
	if tr != nil {
		tr = terr.Add("InitState", errors.New("Problem initilizing services"))
		return tr
	}
	return nil
}

// internal methods

func initServices(dev bool, verbosity int) *terr.Trace {
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
		// init service state
		if name == "http" {
			tr := httpState.InitState(dev, verbosity)
			if tr != nil {
				tr = terr.Add("initServices", errors.New("Unable to initialize http service"))
				return tr
			}
			if Verbosity > 0 {
				terr.Ok("Http service initialized")
			}
		}
	}
	return nil
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
