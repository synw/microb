package types

import (
	"github.com/synw/centcom"
	"github.com/synw/terr"
	"time"
)

type WsServer struct {
	Name       string
	Addr       string
	Key        string
	CmdChanIn  string
	CmdChanOut string
}

type Service struct {
	Name string
	Cmds map[string]*Cmd
	Init func(bool, bool) error
}

type State struct {
	WsServer *WsServer
	Cli      *centcom.Cli
	Services map[string]*Service
	Cmds     map[string]*Cmd
	Conf     *Conf
}

type Cmd struct {
	Id           string
	Name         string
	Date         time.Time
	Service      string
	Args         []interface{}
	From         string
	Status       string
	ErrMsg       string
	Trace        *terr.Trace
	ReturnValues []interface{}
	Exec         interface{}
	ExecCli      interface{}
	ExecAfter    interface{}
}

type Event struct {
	Id      string
	Class   string
	Date    time.Time
	Msg     string
	Service string
	Cmd     *Cmd
	Trace   *terr.Trace
	LogLvl  string
	Data    map[string]interface{}
}

type Conf struct {
	Addr      string
	Key       string
	Name      string
	Services  []string
	RedisAddr string
	RedisDb   int
}
