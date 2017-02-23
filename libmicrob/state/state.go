package state

import (
	"fmt"
	"errors"
	"time"
	"strconv"
	"net/http"
	"github.com/acmacalister/skittles"
	"github.com/synw/microb/libmicrob/datatypes"
	"github.com/synw/microb/libmicrob/conf"
)


var config = conf.GetConf("default")
var Server = &datatypes.Server{Runing:false}
var Verbosity int = 1
var Debug bool = false
var DevMode bool = false
var ListenWs bool = false
var Routes []string
var DbIsOk bool

func InitState(dev_mode bool) error {
	setDevMode(dev_mode)
	err := initState()
	if err != nil {
		return err
	}
	setDebug()
	setVerbosity()
	if Debug == true {
		go func() {
			time.Sleep(1*time.Second)
			msg := "State:\n"+FormatState()
			fmt.Println(msg)
			msg = "Server:\n"+Server.Format()
			fmt.Println(msg)
		}()
	}
	return nil
}

func SetRoutes(routes []string) {
	Routes = routes
}

func FormatState() string {
	var msg string
	d := "off"
	if Debug == true {
		d = "on"
	}
	db := skittles.Red("down")
	if DbIsOk == true {
		db = "up"
	}
	up := skittles.Red("down")
	if Server.Runing == true {
		up = "up"
	}
	cc := "down"
	if ListenWs == true {
		cc = "up"
	}
	msg = msg+" - Server is "+up+"\n"
	msg = msg+" - Database is "+db+"\n"
	msg = msg+" - Commands channel is "+cc+"\n"
	msg = msg+" - Verbosity is set to "+strconv.Itoa(Verbosity)+"\n"
	msg = msg+" - Debug is "+d+"\n"
	if DevMode == true {
		msg = msg+" - Dev mode is on"
	}
	return msg
}

func initState() error {
	err := setServer()
	if err != nil {
		msg := "state.init(): "+err.Error()
		e := errors.New(msg)
		return e
	}
	err = setDefaultDbs()
	if err != nil {
		msg := "state.init(): "+err.Error()
		e := errors.New(msg)
		return e
	}
	return nil
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
		errors.New("state.setDefaultDbs: no database found for role pages")
		return err
	} else {
		Server.PagesDb = pdb
	}
	hdb, err := getDefaultDb("hits")
	if err != nil {
		return err
	}
	if hdb == nil {
		errors.New("state.setDefaultDbs: no database found for role hits")
		return err
	} else {
		Server.HitsDb = hdb
	}
	cdb, err := getDefaultDb("commands")
	if err != nil {
		return err
	}
	if cdb == nil {
		errors.New("state.setDefaultDbs: no database found for role commands")
		return err
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

func setVerbosity() {
	Verbosity = conf.Verbosity
}

func getDefaultDb(role string) (*datatypes.Database, error) {
	dbs, err := getDbsFromConf()
	db := &datatypes.Database{}
	if err != nil {
		e := "state.getDefaultDb: "+err.Error()
		errors.New(e)
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
	e := "state.getDefaultDb: "+msg
	err = errors.New(e)
	return db, err
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
		msg := "state.newDbFromConf: database "+name+" not found in config"
		err := errors.New(msg)
		return db, err
	}
	return db, nil
}
