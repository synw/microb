package types

import (
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
	Name     string
	Cmds     []string
	Init     func(bool, int, bool) *terr.Trace
	Dispatch func(*Cmd) *Cmd
}

type Cmd struct {
	Id           string
	Name         string
	Date         time.Time
	Service      *Service
	Args         []interface{}
	From         string
	Status       string
	ErrMsg       string
	Trace        *terr.Trace
	ReturnValues []interface{}
	Exec         interface{}
}

type Event struct {
	Id      string
	Class   string
	Date    time.Time
	Msg     string
	Service *Service
	Cmd     *Cmd
	Trace   *terr.Trace
	Data    map[string]interface{}
}

type Conf struct {
	Addr     string
	Key      string
	Name     string
	Services []string
}
