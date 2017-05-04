package datatypes

import (
	"time"
	"net/http"
	"github.com/synw/terr"
)


type Server struct {
	Domain string
	HttpHost string
	HttpPort int
	WsHost string
	WsPort int
	WsKey string
	CmdChanIn string
	CmdChanOut string
}

type HttpServer struct {
	Server *Server
	Instance *http.Server
	Running bool
}

type Event struct {
	Class string
	From string
	Message string
	Data map[string]interface{}
	Trace *terr.Trace
}

type Command struct {
	Id string
	Name string
	From string
	Reason string
	Date time.Time
	Args []interface{}
	Status string
	Trace *terr.Trace
	ErrMsg string
	ReturnValues []interface{}
}

type Database struct {
	Type string
	Name string
	Addr string
	User string
	Password string
	Roles []string
	Running bool
}
