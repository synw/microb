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
	crs := Conf["cors"].([]string)
	var cors string
	max := len(crs) + 1
	for i, c := range crs {
		cors = cors + c
		if i < max {
			cors = cors + ","
		}
	}
	host := Conf["host"].(string)
	port := Conf["port"].(int)
	domain := Conf["domain"].(string)
	instance := &http.Server{}
	runing := false
	HttpServer = &datatypes.HttpServer{domain, host, port, instance, runing, cors}
	return nil
}
