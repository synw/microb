package conf

import (
	"fmt"
	"github.com/spf13/viper"
)


func GetConf() map[string]interface{} {
	// set some defaults for conf
	viper.SetConfigName("config")
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
	    panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	conf := make(map[string]interface{})
	conf["domain"] = viper.Get("domain")
	conf["http_host"] = viper.Get("http_host")
	conf["http_port"] = viper.Get("http_port")
	conf["centrifugo_host"] = viper.Get("centrifugo_host")
	conf["centrifugo_port"] = viper.Get("centrifugo_port")
	conf["centrifugo_secret_key"] = viper.Get("centrifugo_secret_key")
	conf["databases"] = viper.Get("databases")
	conf["default_database"] = viper.Get("databases.main")
	conf["staticfiles_host"] = viper.Get("staticfiles_host")
	conf["hits_log"] = viper.Get("hits_log")
	conf["hits_monitor"] = viper.Get("hits_monitor")
	conf["hits_channels"] = viper.Get("hits_channels")
	conf["verbosity"] = viper.Get("verbosity")
	conf["commands_brokers"] = viper.Get("commands_brokers")
	conf["debug"] = viper.Get("debug")
	return conf
}
