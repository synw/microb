package conf

import (
	"errors"
	"github.com/spf13/viper"
	"github.com/synw/terr"
	"github.com/synw/microb/libmicrob/datatypes"
)


func GetConf(name string) (*datatypes.Conf, *terr.Trace) {
	conf, trace := getConf(name)
	if trace != nil {
		trace = terr.Pass("conf.Conf", trace)
		var cf *datatypes.Conf
		return cf, trace
	}
	host := conf["http_host"].(string)
	port := conf["http_port"].(string)
	wshost := conf["centrifugo_host"].(string)	
	wsport := conf["centrifugo_port"].(string)	
	wskey := conf["centrifugo_key"].(string)
	cf := &datatypes.Conf{host, port, wshost, wsport, wskey, 1}
	return cf, nil
}

func getConf(name string) (map[string]interface{},*terr.Trace) {
	// set some defaults for conf
	if name == "dev" {
		viper.SetConfigName("dev_config")
	} else {
		viper.SetConfigName("config")
	}
	viper.AddConfigPath(".")
	viper.SetDefault("http_host", "")
	viper.SetDefault("http_port", "8080")
	viper.SetDefault("centrifugo_host", "localhost")
	viper.SetDefault("centrifugo_port", 8001)
	viper.SetDefault("hits_log", true)
	viper.SetDefault("hits_monitor", true)
	hits_channels := []string{"$microb_hits"}
	viper.SetDefault("hits_channels", hits_channels)
	viper.SetDefault("verbosity", 1)
	viper.SetDefault("debug", false)
	brockers := [0]string{}
	viper.SetDefault("commands_brokers", brockers)
	// get the actual conf
	err := viper.ReadInConfig()
	if err != nil {
		var conf map[string]interface{}
		switch err.(type) {
		case viper.ConfigParseError:
			trace := terr.New("conf.getConf", err)
			return conf, trace
		default:
			err := errors.New("Unable to locate config file")
			trace := terr.New("conf.getConf", err)
			return conf, trace
		}
	}
	conf := make(map[string]interface{})
	conf["domain"] = viper.Get("domain")
	conf["http_host"] = viper.Get("http_host")
	conf["http_port"] = viper.Get("http_port")
	conf["centrifugo_host"] = viper.Get("centrifugo_host")
	conf["centrifugo_port"] = viper.Get("centrifugo_port")
	conf["centrifugo_key"] = viper.Get("centrifugo_key")
	conf["databases"] = viper.Get("databases")
	conf["default_database"] = viper.Get("databases.main")
	conf["staticfiles_host"] = viper.Get("staticfiles_host")
	conf["hits_log"] = viper.Get("hits_log")
	conf["hits_monitor"] = viper.Get("hits_monitor")
	conf["hits_channels"] = viper.Get("hits_channels")
	/*conf["verbosity"] = viper.Get("verbosity")
	v := conf["verbosity"]
	Verbosity, _ = v.(int)*/
	conf["commands_brokers"] = viper.Get("commands_brokers")
	return conf, nil
}
