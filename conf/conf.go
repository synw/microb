package conf

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/synw/microb/db/datatypes"
)


func GetConf() map[string]interface{} {
	// set some defaults for conf
	viper.SetConfigName("microb_dev_config")
	viper.AddConfigPath(".")
	viper.SetDefault("http_host", ":8080")
	viper.SetDefault("centrifugo_host", "localhost")
	viper.SetDefault("centrifugo_port", 8001)
	viper.SetDefault("hits_log", true)
	viper.SetDefault("hits_monitor", true)
	hits_channels := []string{"$microb_hits"}
	viper.SetDefault("hits_channels", hits_channels)
	viper.SetDefault("commands_brokers", []string{"default"})
	// get the actual conf
	err := viper.ReadInConfig()
	if err != nil {
	    panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	conf := make(map[string]interface{})
	conf["http_host"] = viper.Get("http_host")
	conf["centrifugo_host"] = viper.Get("centrifugo_host")
	conf["centrifugo_port"] = viper.Get("centrifugo_port")
	conf["centrifugo_secret_key"] = viper.Get("centrifugo_secret_key")
	conf["databases"] = viper.Get("databases").(interface{})
	conf["main_database"] = viper.Get("databases.main")
	conf["domain"] = viper.Get("domain")
	conf["hits_log"] = viper.Get("hits_log")
	conf["hits_monitor"] = viper.Get("hits_monitor")
	conf["hits_channels"] = viper.Get("hits_channels")
	conf["commands_brokers"] = viper.Get("commands_brokers")
	return conf
}

var Config = GetConf()

func GetMainDb() *datatypes.Database {
	db := make(map[string]string)
	main_db := Config["main_database"].(map[string]interface{})
	db["type"] = main_db["type"].(string)
	db["host"] = main_db["host"].(string)
	db["port"] = main_db["port"].(string)
	db["user"] = main_db["user"].(string)
	db["password"] = main_db["password"].(string)
	database := &datatypes.Database{db["type"], "main", db["host"], db["port"], db["user"], db["password"]}
	return database
}

func commandsTransports() []string {
	var ts []string
	cts := Config["commands_brokers"].([]string)
	db_type := GetMainDb().Type
	// check for defaults
	is_default := false
	for _, transp := range cts {
		if transp == "default" {
			is_default = true
			if db_type == "rethinkdb" {
				ts = []string{"changefeeds"}
			}
		}
	}
	if is_default == false {
		ts = cts
	}
	return ts
}

func ListenToChangefeeds() bool {
	listen := false
	transports := commandsTransports()
	for _, val := range transports {
		if val == "changefeeds" {
			listen = true
			break
		}
	}
	return listen
}
