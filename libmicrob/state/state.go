package state

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/synw/centcom"
	"github.com/synw/microb/libmicrob/conf"
	"github.com/synw/microb/libmicrob/types"
	"github.com/synw/terr"
)

var Debug = true
var Verbosity int = 1
var Cli *centcom.Cli
var Server = &types.Server{}
var Conf = &types.Conf{}
var Dev bool
var Services = make(map[string]*types.Service)
var Logger *logrus.Logger

func InitState(dev bool, verbosity int) (*types.Conf, *terr.Trace) {
	Verbosity = verbosity
	Dev = dev
	// conf
	Conf, tr := conf.GetConf(dev)
	if tr != nil {
		return Conf, tr
	}
	// command channel server
	server, tr := conf.GetServer(Conf)
	if tr != nil {
		tr = terr.Pass("stateInit.State", tr)
		return Conf, tr
	}
	Server = server
	// Websockets client
	if Verbosity > 0 {
		fmt.Println("Initializing commands transport layer ...")
	}
	cli, tr := initWsCli()
	if tr != nil {
		tr = terr.Pass("state.InitState", tr)
		return Conf, tr
	}
	Cli = cli
	return Conf, nil
}

// internal methods

func initWsCli() (*centcom.Cli, *terr.Trace) {
	cli := centcom.NewClient(Server.WsAddr, Server.WsKey)
	err := centcom.Connect(cli)
	if err != nil {
		trace := terr.New("initWsCli", err)
		var cli *centcom.Cli
		return cli, trace
	}
	cli.IsConnected = true
	if Verbosity > 1 {
		terr.Ok("Websockets client connected")
	}
	err = cli.CheckHttp()
	if err != nil {
		trace := terr.New("InitCli", err)
		return cli, trace
	}
	if Verbosity > 1 {
		terr.Ok("Websockets http transport ready")
	}
	return cli, nil
}
