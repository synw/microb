package conf

import (
	"fmt"
	"github.com/spf13/viper"
)


func GetConf() map[string]interface{} {
	viper.SetConfigName("microb_dev_config_pg")
	viper.AddConfigPath(".")
	viper.SetDefault("http_host", ":8080")
	viper.SetDefault("centrifugo_host", "localhost")
	viper.SetDefault("centrifugo_port", 8001)
	viper.SetDefault("db_type", "rethinkdb")
	viper.SetDefault("db_host", "localhost")
	viper.SetDefault("db_port", 28015)
	viper.SetDefault("db_user", "admin")
	viper.SetDefault("db_password", "")
	viper.SetDefault("hits_log", true)
	viper.SetDefault("hits_monitor", true)
	hits_channels := []string{"$microb_hits"}
	viper.SetDefault("hits_channels", hits_channels)
	viper.SetDefault("commands_transport", []string{"default"})
	err := viper.ReadInConfig()
	if err != nil {
	    panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	conf := make(map[string]interface{})
	conf["http_host"] = viper.Get("http_host")
	conf["centrifugo_host"] = viper.Get("centrifugo_host")
	conf["centrifugo_port"] = viper.Get("centrifugo_port")
	conf["centrifugo_secret_key"] = viper.Get("centrifugo_secret_key")
	conf["db_type"] = viper.Get("db_type")
	conf["db_host"] = viper.Get("db_host")
	conf["db_port"] = viper.Get("db_port")
	conf["db_user"] = viper.Get("db_user")
	conf["db_password"] = viper.Get("db_password")
	conf["domain"] = viper.Get("domain")
	conf["hits_log"] = viper.Get("hits_log")
	conf["hits_monitor"] = viper.Get("hits_monitor")
	conf["hits_channels"] = viper.Get("hits_channels")
	conf["commands_transport"] = viper.Get("commands_transport")
	return conf
}
