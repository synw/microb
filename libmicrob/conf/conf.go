package conf

import (
	"errors"
	"github.com/spf13/viper"
	"github.com/synw/microb/libmicrob/types"
	"github.com/synw/terr"
)

func GetComChan(name string) (string, string) {
	comchan_in := "cmd:$" + name + "_in"
	comchan_out := "cmd:$" + name + "_out"
	return comchan_in, comchan_out
}

func GetServer(conf *types.Conf) (*types.Server, *terr.Trace) {
	comchan_in, comchan_out := GetComChan(conf.Name)
	s := &types.Server{conf.Name, conf.WsHost, conf.WsPort, conf.WsKey, comchan_in, comchan_out}
	return s, nil
}

func GetConf(dev bool) (*types.Conf, *terr.Trace) {
	name := "normal"
	if dev {
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
	viper.SetDefault("centrifugo_host", "localhost")
	viper.SetDefault("centrifugo_port", 8001)
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
		viper.Get("centrifugo_host").(string),
		int(viper.Get("centrifugo_port").(float64)),
		viper.Get("centrifugo_key").(string),
		viper.Get("name").(string),
		services,
	}
	return conf, nil
}
