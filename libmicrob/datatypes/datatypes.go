package datatypes

import (
	"github.com/synw/terr"
)


type Conf struct {
	HttpHost string
	HttpPort string
	WsHost string
	WsPort string
	WkKey string
	Verbosity int
}

type Event struct {
	Class string
	From string
	Message string
	Data map[string]interface{}
	Trace *terr.Trace
}
