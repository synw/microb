package conf

import (
	"errors"
	"github.com/spf13/viper"
	"github.com/synw/microb/libmicrob/datatypes"
	"github.com/synw/terr"
)

func GetServer(dev bool) (*datatypes.Server, *terr.Trace) {
	conf, tr := GetConf(dev)
	if tr != nil {
		s := &datatypes.Server{}
		return s, tr
	}
	name := conf["name"].(string)
	wshost := conf["centrifugo_host"].(string)
	wsport := conf["centrifugo_port"].(int)
	wskey := conf["centrifugo_key"].(string)

	comchan_in := "cmd:$" + name + "_in"
	comchan_out := "cmd:$" + name + "_out"
	s := &datatypes.Server{name, wshost, wsport, wskey, comchan_in, comchan_out}
	return s, nil
}

func GetDefaultDb(name string) (*datatypes.Database, *terr.Trace) {
	conf, tr := getConf(name)
	if tr != nil {
		d := &datatypes.Database{}
		return d, tr
	}
	db := conf["default_database"].(map[string]interface{})
	dbtype := db["type"].(string)
	addr := db["addr"].(string)
	user := db["user"].(string)
	password := db["password"].(string)
	r := db["roles"].([]interface{})
	d := db["dbs"].(map[string]interface{})
	t := db["tables"].(map[string]interface{})
	var roles []string
	var dbs = make(map[string]string)
	dbs["documents"] = ""
	dbs["metrics"] = ""
	dbs["commands"] = ""
	var tables = make(map[string]string)
	tables["documents"] = ""
	tables["metrics"] = ""
	tables["commands"] = ""
	for _, role := range r {
		srole := role.(string)
		roles = append(roles, srole)
		if d[srole] != nil {
			dbs[srole] = d[srole].(string)
		}
		if t[srole] != nil {
			tables[srole] = t[srole].(string)
		}
	}
	dbobj := &datatypes.Database{dbtype, "default", addr, user, password, roles, dbs, tables, false}
	return dbobj, nil
}

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
	viper.SetDefault("centrifugo_host", "localhost")
	viper.SetDefault("centrifugo_port", 8001)
	viper.SetDefault("hits_log", true)
	viper.SetDefault("hits_monitor", true)
	hits_channels := []string{"$microb_hits"}
	viper.SetDefault("hits_channels", hits_channels)
	viper.SetDefault("debug", false)
	viper.SetDefault("name", "localhost")
	viper.SetDefault("services", []string{})
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
	conf["centrifugo_host"] = viper.Get("centrifugo_host").(string)
	conf["centrifugo_port"] = int(viper.Get("centrifugo_port").(float64))
	conf["centrifugo_key"] = viper.Get("centrifugo_key").(string)
	conf["name"] = viper.Get("name").(string)
	conf["services"] = viper.Get("services").([]interface{})
	conf["databases"] = viper.Get("databases")
	conf["default_database"] = viper.Get("databases.default")
	conf["hits_log"] = viper.Get("hits_log")
	conf["hits_monitor"] = viper.Get("hits_monitor")
	conf["hits_channels"] = viper.Get("hits_channels")
	return conf, nil
}
