package conf

import (
	"errors"
	"github.com/spf13/viper"
	"github.com/synw/terr"
	"github.com/synw/microb/libmicrob/datatypes"
)


func GetServer(name string) (*datatypes.Server, *terr.Trace) {
	conf, tr := getConf(name)
	if tr != nil {
		s := &datatypes.Server{}
		return s, tr
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

func GetDefaultDb(name string) (*datatypes.Database, *terr.Trace) {
	conf, tr := getConf(name)
	if tr != nil {
		d := &datatypes.Database{}
		return d, tr
	}
	db := conf["default_database"].(map[string]interface {})
	dbtype := db["type"].(string)
	addr := db["addr"].(string)
	user := db["user"].(string)	
	password := db["password"].(string)
	r := db["roles"].([]interface{})
	var roles []string
	for _, role := range(r) {
		roles = append(roles, role.(string))
	}
	d := &datatypes.Database{dbtype, "default", addr, user, password, roles, false}
	return d, nil
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
	conf["default_database"] = viper.Get("databases.default")
	conf["staticfiles_host"] = viper.Get("staticfiles_host")
	conf["hits_log"] = viper.Get("hits_log")
	conf["hits_monitor"] = viper.Get("hits_monitor")
	conf["hits_channels"] = viper.Get("hits_channels")
	return conf, nil
}
