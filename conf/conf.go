package conf

import (
	"fmt"
	"github.com/spf13/viper"
)


func GetConf() map[string]interface{} {
	viper.SetConfigName("microb_config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./etc/microb")
	viper.AddConfigPath("$HOME/.microb")
	viper.SetDefault("http_port", 8080)
	//viper.SetDefault("centrifugo_host", "localhost")
	//viper.SetDefault("centrifugo_port", 8001)
	viper.SetDefault("db_type", "rethinkdb")
	viper.SetDefault("db_host", "localhost")
	viper.SetDefault("db_port", 28015)
	viper.SetDefault("db_user", "admin")
	viper.SetDefault("db_password", "")
	viper.SetDefault("database", "microb")
	viper.SetDefault("table", "pages")
	viper.SetDefault("domain", "localhost")
	err := viper.ReadInConfig()
	if err != nil {
	    panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	conf := make(map[string]interface{})
	conf["http_port"] = viper.Get("http_port")
	//conf["centrifugo_host"] = viper.Get("centrifugo_host")
	//conf["centrifugo_port"] = viper.Get("centrifugo_port")
	//conf["centrifugo_secret_key"] = viper.Get("centrifugo_secret_key")
	conf["db_type"] = viper.Get("db_type")
	conf["db_host"] = viper.Get("db_host")
	conf["db_port"] = viper.Get("db_port")
	conf["db_user"] = viper.Get("db_user")
	conf["db_password"] = viper.Get("db_password")
	conf["db_type"] = viper.Get("db_type")
	conf["database"] = viper.Get("database")
	conf["table"] = viper.Get("table")
	conf["domain"] = viper.Get("domain")
	return conf
}
