package state

import (
	//"fmt"
	"errors"
	"time"
	"strconv"
	"net/http"
	"github.com/synw/microb/libmicrob/datatypes"
	"github.com/synw/microb/libmicrob/events"
	"github.com/synw/microb/libmicrob/conf"
)


var config = conf.GetConf("default")
var Server = &datatypes.Server{Runing:false}
var Verbosity int = getVerbosity()
var Debug bool = false
var DevMode bool = false
var ListenWs bool = false
var Routes []string
var DbIsOk bool

func InitState(dev_mode bool) {
	setDevMode(dev_mode)
	initState()
	setDebug()
	if Debug == true {
		go func() {
			time.Sleep(1*time.Second)
			events.Debug("State:", "\n", printState())
			events.Debug("Server:", "\n", Server.Format())
		}()
	}
}

func SetRoutes(routes []string) {
	Routes = routes
}

func initState() {
	err := setServer()
	if err != nil {
		events.Error("state.init()", err)
		return
	}
	err = setDefaultDbs()
	if err != nil {
		events.Error("state.init()", err)
		return
	}
}

func setServer() (error) {
	domain := config["domain"].(string)
	http_host :=  config["http_host"].(string)
	http_port := config["http_port"].(string)
	websockets_host := config["centrifugo_host"].(string)
	websockets_port := config["centrifugo_port"].(string)
	websockets_key := config["centrifugo_secret_key"].(string)
	pages_db := Server.PagesDb
	hits_db := Server.HitsDb
	commands_db := Server.CommandsDb
	var srv *http.Server
	server := &datatypes.Server{domain, http_host, http_port, websockets_host, websockets_port, websockets_key, pages_db, hits_db, commands_db, false, srv}
	Server = server
	return nil
}

func setDefaultDbs() error {
	pdb, err := getDefaultDb("pages")
	if err != nil {
		return err
	}
	if pdb == nil {
		events.ErrMsg("state.setDefaultDbs", "No database found for role pages")
	} else {
		Server.PagesDb = pdb
	}
	hdb, err := getDefaultDb("hits")
	if err != nil {
		return err
	}
	if hdb == nil {
		events.ErrMsg("state.setDefaultDbs", "No database found for role hits")
	} else {
		Server.HitsDb = hdb
	}
	cdb, err := getDefaultDb("commands")
	if err != nil {
		return err
	}
	if cdb == nil {
		events.ErrMsg("state.setDefaultDbs", "No database found for role commands")
	} else {
		Server.CommandsDb = cdb
	}
	return nil
}

func setDevMode(dev_mode bool) {
	if dev_mode == true {
		config = conf.GetConf("dev")
		DevMode = true
	} else {
		DevMode = false
	}
}

func setDebug() {
	d := config["debug"].(bool)
	Debug = d
}

func getDefaultDb(role string) (*datatypes.Database, error) {
	dbs, err := getDbsFromConf()
	db := &datatypes.Database{}
	if err != nil {
		events.Error("state.getDefaultDb", err)
		return db, err
	}
	for _, cdb := range dbs {
		for _, r := range cdb.Roles {
			if r == role {
				db = cdb
				return db, nil
			}
		}
	}
	msg := "Database not found for role "+role
	events.ErrMsg("state.getDefaultDb", msg)
	return db, nil
}

func getDbsFromConf() (map[string]*datatypes.Database, error) {
	dbs := make(map[string]*datatypes.Database)
	dbs_conf := config["databases"].(map[string]interface{})
	for db_name, _ := range(dbs_conf) {
		db, err := newDbFromConf(db_name)
		if err != nil {
			return dbs, err
		}
		dbs[db_name] = db
	}
	return dbs, nil
}

func newDbFromConf(name string) (*datatypes.Database, error) {
	dbs_conf := config["databases"].(map[string]interface{})
	var db *datatypes.Database
	// grab the db in the config
	for db_name, db_vals := range dbs_conf {
		db_vals := db_vals.(map[string]interface{})
		if db_name == name {
			dbtype := db_vals["type"].(string)
			host := db_vals["host"].(string)
			port := db_vals["port"].(string)
			user := db_vals["user"].(string)
			password := db_vals["password"].(string)
			roles_i := db_vals["roles"].([]interface{})
			var roles []string
			for _, r := range roles_i {
				roles = append(roles, r.(string))
			}
			db = &datatypes.Database{dbtype, db_name, host, port, user, password, roles}
			break
		}
	}
	if db == nil {
		msg := "Database "+name+" not found in config"
		err := errors.New(msg)
		events.Error("state.newDbFromConf", err)
		return db, err
	}
	return db, nil
}

func getVerbosity() int {
	v := config["verbosity"].(int)
	return v
}

func printState() string {
	dm := "off"
	if DevMode == true {
		dm = "on"
	}
	msg := " - Dev mode is "+dm+"\n"
	d := "off"
	if Debug == true {
		d = "on"
	}
	msg = msg+" - Debug is "+d+"\n"
	msg = msg+" - Verbosity is "+strconv.Itoa(Verbosity)
	return msg
}
