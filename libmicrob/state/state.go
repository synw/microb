package state

import (
	"fmt"
	"github.com/synw/terr"
	"github.com/synw/centcom"
	"github.com/synw/centcom/ws"
	"github.com/synw/microb/libmicrob/conf"
	"github.com/synw/microb/libmicrob/datatypes"
)


var Server *datatypes.Server = &datatypes.Server{}
var Verbosity int = 1
var Cli *ws.Cli


func InitState(is_dev bool, verbosity int) (*ws.Cli, *terr.Trace) {
	Verbosity = verbosity	
	name := "normal"
	if is_dev == true {
		name = "dev"
		fmt.Println("Dev mode is on")
	}
	server, trace := conf.GetServer(name)
	if trace != nil {
		trace = terr.Pass("stateInit.State", trace)
		var cli *ws.Cli
		return cli, trace
	}
	Server = server
	cli, trace := initWsCli()
	if trace != nil {	
		trace = terr.Pass("state.InitState", trace)
		return cli, trace
	}
	return cli, nil
}

// internal methods

func initWsCli() (*ws.Cli, *terr.Trace) {
	cli := ws.NewClient(Server.WsHost, Server.WsPort, Server.WsKey)
	cli, err := centcom.Connect(cli)
	if err != nil {
		trace := terr.New("ws.InitCli", err)
		var cli *ws.Cli
		return cli, trace
	}
	cli.IsConnected = true
	if Verbosity > 1 {
		fmt.Println(terr.Ok("Websockets client connected"))
	}
	cli, err = cli.CheckHttp()
	if err != nil {
		trace := terr.New("ws.InitCli", err)
		return cli, trace
	}
	if Verbosity > 1 {
		fmt.Println(terr.Ok("Websockets http transport ready"))
	}	
	return cli, nil
}
