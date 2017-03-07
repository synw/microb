package datatypes

import (
	"time"
	"github.com/synw/terr"
)


type Server struct {
	Domain string
	HttpHost string
	HttpPort int
	WsHost string
	WsPort int
	WsKey string
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
	Error error
	ReturnValues []interface{}
}
