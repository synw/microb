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
	viper.SetDefault("host", "")
	viper.SetDefault("domain", "localhost")
	viper.SetDefault("port", 8080)
	viper.SetDefault("cors", []string{})
	viper.SetDefault("staticfiles_host", "")
	viper.SetDefault("staticfiles_port", 3000)
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
	comchan := "$" + conf["domain"].(string)
	viper.SetDefault("command_channel", comchan)
	conf["host"] = viper.Get("host").(string)
	conf["domain"] = viper.Get("domain").(string)
	conf["port"] = viper.Get("port").(int)
	conf["staticfiles_host"] = viper.Get("staticfiles_host").(string)
	conf["staticfiles_port"] = viper.Get("staticfiles_port").(int)
	conf["cors"] = viper.Get("cors").([]string)
	return conf, nil
}
