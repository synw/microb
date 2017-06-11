package datatypes

import (
	"github.com/synw/terr"
	"time"
)

type Conf struct {
	WsHost   string
	WsPort   int
	WsKey    string
	Name     string
	Services []string
}

type Server struct {
	Name       string
	WsHost     string
	WsPort     int
	WsKey      string
	CmdChanIn  string
	CmdChanOut string
}

type Event struct {
	Class   string
	From    string
	Service string
	Date    time.Time
	Msg     string
	Err     error
	Data    map[string]interface{}
}

type Service struct {
	Name     string
	Cmds     []string
	Requires []*Service
}

type Command struct {
	Id           string
	Service      string
	Name         string
	From         string
	Reason       string
	Date         time.Time
	Args         []interface{}
	Status       string
	Trace        *terr.Trace
	ErrMsg       string
	ReturnValues []interface{}
}
