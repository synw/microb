package state

import (
	"errors"
	"github.com/synw/microb/libmicrob/datatypes"
	"github.com/synw/microb/libmicrob/events"
	"github.com/synw/microb/libmicrob/conf"
)

var Config = conf.GetConf()
var Server = &datatypes.Server{Runing:false}
var Verbosity int = getVerbosity()
var Debug bool = getDebug()

func init() {
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

func setServer() ( error) {
	domain := Config["domain"].(string)
	http_host :=  Config["http_host"].(string)
	http_port := Config["http_port"].(string)
	websockets_host := Config["centrifugo_host"].(string)
	websockets_port := Config["centrifugo_port"].(string)
	websockets_key := Config["centrifugo_secret_key"].(string)
	pages_db := Server.PagesDb
	hits_db := Server.HitsDb
	commands_db := Server.CommandsDb
	server := &datatypes.Server{domain, http_host, http_port, websockets_host, websockets_port, websockets_key, pages_db, hits_db, commands_db, false}
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
/*
func getDefaultDb(role string) (*datatypes.Database, error) {
	var db *datatypes.Database
	default_roles, err := getRolesFromConf()
	if err != nil {
		return db, err
	} else {
		events.Error("state.getDefaultDb", err)
	}
	db = default_roles[role]
	return db, nil
}*/
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
/*
func getDbRoles() (map[string]*datatypes.Database, error) {
	var rolemap map[string]*datatypes.Database
	dbs, err := getDbsFromConf()
	if err != nil {
		return rolemap, err
	}
	roles, err := getRolesFromConf()
	if err != nil {
		return rolemap, err
	}
	for _, db := range dbs {
		for role, cdb := range roles {
			if cdb == db {
				rolemap[role] = cdb
			}
		}
	}
	return rolemap, err
}*/
/*
func getRolesFromConf() (map[string]*datatypes.Database, error) {
	conf_roles := Config["database_roles"].(map[string]interface{})
	roles := make(map[string]*datatypes.Database)
	for role, db_name := range(conf_roles) {
		r, err := newDbFromConf(db_name.(string))
		if err != nil {
			events.Error("state.getRolesFromConf", err)
			return roles, err
		}
		roles[role] = r
	}
	return roles, nil
}*/

func getDbsFromConf() (map[string]*datatypes.Database, error) {
	dbs := make(map[string]*datatypes.Database)
	dbs_conf := Config["databases"].(map[string]interface{})
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
	dbs_conf := Config["databases"].(map[string]interface{})
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
	v := Config["verbosity"].(int)
	return v
}

func getDebug() bool {
	d := Config["debug"].(bool)
	return d
}
