package state

import (
	"fmt"
	"github.com/synw/terr"
	"github.com/synw/centcom"
	"github.com/synw/microb/libmicrob/conf"
	"github.com/synw/microb/libmicrob/datatypes"
)


var Debug = true
var Server = &datatypes.Server{}
var Verbosity int = 1
var Cli *centcom.Cli
var HttpServer = datatypes.HttpServer{}
var DocDb = &datatypes.Database{}

func InitState(name string, verbosity int) *terr.Trace {
	Verbosity = verbosity
	// microb server
	server, trace := conf.GetServer(name)
	if trace != nil {
		trace = terr.Pass("stateInit.State", trace)
		return trace
	}
	Server = server
	// Websockets client
	if Verbosity > 0 {
		fmt.Println("Initializing commands transport layer ...")
	}
	cli, trace := initWsCli()
	if trace != nil {	
		trace = terr.Pass("state.InitState", trace)
		return trace
	}
	Cli = cli
	return nil
}

// internal methods

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
