package metadata

import (
	"github.com/synw/microb/libmicro/conf"
    "github.com/synw/microb/libmicro/datatypes"
)


var Config = conf.GetConf()

func GetVerbosity() int {
	v := Config["verbosity"].(int)
	return v
}

func GetServer() *datatypes.Server {
	server := &datatypes.Server{Config["domain"].(string), Config["http_host"].(string), Config["http_port"].(string)}
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
