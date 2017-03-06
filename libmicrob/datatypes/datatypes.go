package datatypes

import (
	"github.com/synw/terr"
)


type Conf struct {
	HttpHost string
	HttpPort int
	WsHost string
	WsPort int
	WsKey string
	Verbosity int
}

type Event struct {
	Class string
	From string
	Message string
	Data map[string]interface{}
	Trace *terr.Trace
}
