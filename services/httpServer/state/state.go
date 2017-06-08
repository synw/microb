package state

import (
	"github.com/synw/microb/services/httpServer/conf"
	"github.com/synw/microb/services/httpServer/datatypes"
	"github.com/synw/terr"
	"net/http"
)

var HttpServer = &datatypes.HttpServer{}
var Conf map[string]interface{}

func InitState(dev bool, verbosity int) *terr.Trace {
	Conf, tr := conf.GetConf(dev)
	if tr != nil {
		return tr
	}
	crs := Conf["http_cors"].([]string)
	var cors string
	max := len(crs) + 1
	for i, c := range crs {
		cors = cors + c
		if i < max {
			cors = cors + ","
		}
	}
	host := Conf["http_host"].(string)
	port := Conf["http_port"].(int)
	domain := Conf["http_domain"].(string)
	instance := &http.Server{}
	runing := false
	HttpServer = &datatypes.HttpServer{domain, host, port, instance, runing, cors}
	return nil
}
