package conf

import (
	"errors"
	"github.com/spf13/viper"
	"github.com/synw/terr"
)

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
	viper.SetDefault("http_host", "")
	viper.SetDefault("http_domain", "localhost")
	viper.SetDefault("http_port", 8080)
	viper.SetDefault("http_cors", []interface{}{})
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
	conf["http_domain"] = viper.GetString("http_domain")
	conf["http_host"] = viper.GetString("http_host")
	conf["http_port"] = viper.GetInt("http_port")
	conf["http_cors"] = viper.GetStringSlice("http_cors")
	return conf, nil
}
