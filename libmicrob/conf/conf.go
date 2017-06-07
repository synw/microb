package conf

import (
	"errors"
	"github.com/spf13/viper"
	"github.com/synw/microb/libmicrob/datatypes"
	"github.com/synw/terr"
)

func GetComChan(name string) (string, string) {
	comchan_in := "cmd:$" + name + "_in"
	comchan_out := "cmd:$" + name + "_out"
	return comchan_in, comchan_out
}

func GetServer(dev bool) (*datatypes.Server, *terr.Trace) {
	conf, tr := GetConf(dev)
	if tr != nil {
		s := &datatypes.Server{}
		return s, tr
	}
	name := conf["name"].(string)
	wshost := conf["centrifugo_host"].(string)
	wsport := conf["centrifugo_port"].(int)
	wskey := conf["centrifugo_key"].(string)
	comchan_in, comchan_out := GetComChan(name)
	s := &datatypes.Server{name, wshost, wsport, wskey, comchan_in, comchan_out}
	return s, nil
}

func GetConf(dev bool) (map[string]interface{}, *terr.Trace) {
	name := "normal"
	if dev {
		name = "dev"
	}
	return getConf(name)
}

func getConf(name string) (map[string]interface{}, *terr.Trace) {
	// set some defaults for conf
	if name == "dev" {
		viper.SetConfigName("dev_config")
	} else {
		viper.SetConfigName("config")
	}
	viper.AddConfigPath(".")
	viper.SetDefault("centrifugo_host", "localhost")
	viper.SetDefault("centrifugo_port", 8001)
	viper.SetDefault("debug", false)
	viper.SetDefault("name", "localhost")
	viper.SetDefault("services", []string{})
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
	conf["centrifugo_host"] = viper.Get("centrifugo_host").(string)
	conf["centrifugo_port"] = int(viper.Get("centrifugo_port").(float64))
	conf["centrifugo_key"] = viper.Get("centrifugo_key").(string)
	conf["name"] = viper.Get("name").(string)
	conf["services"] = viper.Get("services").([]interface{})
	return conf, nil
}
