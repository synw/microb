package conf

import (
	"errors"
	"github.com/spf13/viper"
	"github.com/synw/microb/libmicrob/types"
	"github.com/synw/terr"
)

func GetConf(dev string) (*types.Conf, *terr.Trace) {
	name := "normal"
	if dev != "zero" {
		name = "dev"
	}
	return getConf(name)
}

func getConf(name string) (*types.Conf, *terr.Trace) {
	// set some defaults for conf
	if name == "dev" {
		viper.SetConfigName("dev_config")
	} else {
		viper.SetConfigName("config")
	}
	viper.AddConfigPath(".")
	viper.SetDefault("centrifugo_addr", "localhost:8001")
	viper.SetDefault("debug", false)
	viper.SetDefault("name", "localhost")
	viper.SetDefault("services", []string{})
	// get the actual conf
	err := viper.ReadInConfig()
	if err != nil {
		conf := &types.Conf{}
		switch err.(type) {
		case viper.ConfigParseError:
			tr := terr.New("conf.getConf", err)
			return conf, tr
		default:
			err := errors.New("Unable to locate config file")
			tr := terr.New("conf.getConf", err)
			return conf, tr
		}
	}
	services := []string{"info"}
	for _, s := range viper.Get("services").([]interface{}) {
		services = append(services, s.(string))
	}
	conf := &types.Conf{
		viper.Get("centrifugo_addr").(string),
		viper.Get("centrifugo_key").(string),
		viper.Get("name").(string),
		services,
	}
	return conf, nil
}
