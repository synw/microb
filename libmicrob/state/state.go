package state

import (
	"github.com/synw/terr"
	"github.com/synw/microb/libmicrob/conf"
	"github.com/synw/microb/libmicrob/datatypes"
)


var Conf *datatypes.Conf = &datatypes.Conf{}
var Verbosity int = 1


func InitState(is_dev bool, verbosity int) *terr.Trace {
	Verbosity = verbosity	
	cf := "config"
	if is_dev == true {
		cf = "dev_config"
	}
	_, trace := conf.GetConf(cf)
	if trace != nil {
		trace = terr.Pass("state.GetConf", trace)
		return trace
	}
	return nil
}
