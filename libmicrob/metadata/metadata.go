package metadata

import (
	//"fmt"
	"github.com/synw/microb/libmicrob/conf"
    "github.com/synw/microb/libmicrob/datatypes"
)


var Config = conf.GetConf()

func GetConfString(param string) string {
	domain := Config[param].(string)
	return domain
}

func GetVerbosity() int {
	v := Config["verbosity"].(int)
	return v
}

func IsDebug() bool {
	d := Config["debug"].(bool)
	return d
}

func GetServer() *datatypes.Server {
	domain := Config["domain"].(string)
	http_host :=  Config["http_host"].(string)
	http_port := Config["http_port"].(string)
	websockets_host := Config["centrifugo_host"].(string)
	websockets_port := Config["centrifugo_port"].(string)
	websockets_key := Config["centrifugo_secret_key"].(string)
	server := &datatypes.Server{domain, http_host, http_port, websockets_host, websockets_port, websockets_key}
	return server
}

func GetDefaultDatabase() *datatypes.Database {
	db := make(map[string]string)
	main_db := Config["default_database"].(map[string]interface{})
	db["type"] = main_db["type"].(string)
	db["host"] = main_db["host"].(string)
	db["port"] = main_db["port"].(string)
	db["user"] = main_db["user"].(string)
	db["password"] = main_db["password"].(string)
	database := &datatypes.Database{db["type"], "main", db["host"], db["port"], db["user"], db["password"]}
	return database
}

func GetMainDatabase() *datatypes.Database {
	return GetDefaultDatabase()
}

func GetCommandsBrockers() []string {
	conf_brokers := Config["commands_brokers"].([]interface{})
	nb := len(conf_brokers)
	brokers  := make([]string, nb)
	for i, broker := range(conf_brokers) {
		brokers[i] = broker.(string)
	}
	if len(brokers) == 0 {
		b := []string{}
		return b
	}
	return brokers
}

func IsWebsocketsBrocker() bool {
	brockers := GetCommandsBrockers()
	if brockers != nil {
		for _, b := range(brockers) {
			if b == "websockets" {
				return true
			}
		}
	}
	return false
}
