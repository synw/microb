package conf

import (
	"errors"
	"github.com/spf13/viper"
	"github.com/synw/microb/libmicrob/types"
	"github.com/synw/terr"
	"os"
	"path/filepath"
)

func getBasePath() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	cp := filepath.Dir(ex)
	return cp
}

func getComChan(name string) (string, string) {
	comchan_in := "cmd:$" + name + "_in"
	comchan_out := "cmd:$" + name + "_out"
	return comchan_in, comchan_out
}

func GetServer(conf *types.Conf) (*types.WsServer, *terr.Trace) {
	comchan_in, comchan_out := getComChan(conf.Name)
	s := &types.WsServer{conf.Name, conf.Addr, conf.Key, comchan_in, comchan_out}
	return s, nil
}

func GetConf() (*types.Conf, *terr.Trace) {
	// set some defaults for conf
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetDefault("centrifugo_addr", "localhost:8001")
	viper.SetDefault("debug", false)
	viper.SetDefault("name", "localhost")
	viper.SetDefault("services", []string{})
	viper.SetDefault("redis_addr", ":6379")
	viper.SetDefault("redis_db", 0)
	dbpath := getBasePath() + "/logs.sqlite"
	viper.SetDefault("logsDbAddr", dbpath)
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
	var services []string
	for _, s := range viper.Get("services").([]interface{}) {
		services = append(services, s.(string))
	}
	conf := &types.Conf{
		viper.Get("centrifugo_addr").(string),
		viper.Get("centrifugo_key").(string),
		viper.Get("name").(string),
		services,
		viper.Get("redis_addr").(string),
		viper.Get("redis_db").(int),
		viper.Get("logsDbAddr").(string),
	}
	return conf, nil
}
