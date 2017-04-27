package conf

import (
	"errors"
	"github.com/spf13/viper"
	"github.com/synw/terr"
	"github.com/synw/microb/libmicrob/datatypes"
)


func GetServer(name string) (*datatypes.Server, *terr.Trace) {
	conf, trace := getConf(name)
	if trace != nil {
		trace = terr.Pass("conf.GetServer", trace)
		var s *datatypes.Server
		return s, trace
	}
	domain := conf["domain"].(string)
	host := conf["http_host"].(string)
	port := int(conf["http_port"].(float64))
	wshost := conf["centrifugo_host"].(string)	
	wsport := int(conf["centrifugo_port"].(float64))	
	wskey := conf["centrifugo_key"].(string)
	comchan_in := "cmd:$"+domain+"_in"
	comchan_out := "cmd:$"+domain+"_out"
	s := &datatypes.Server{domain, host, port, wshost, wsport, wskey, comchan_in, comchan_out}
	return s, nil
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
	viper.SetDefault("debug", false)
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
	comchan := "$"+conf["domain"].(string)
	viper.SetDefault("command_channel", comchan)
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
	conf["command_channel"] = viper.Get("command_channel")
	return conf, nil
}
