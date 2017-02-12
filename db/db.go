package db

import (
	"sync"
	"github.com/synw/microb/conf"
	"github.com/synw/microb/db/datatypes"
	"github.com/synw/microb/db/rethinkdb"
)

var Config = conf.GetConf()
var Backend = Config["db_type"].(string)

func SaveCommand(command *datatypes.Command, wg *sync.WaitGroup) {
	if Backend == "rethinkdb" {
		rethinkdb.SaveCommand(command, wg)
	}
	return
}

func GetFromDb(url string)  (map[string]interface{}, bool) {
	res := make(map[string]interface{})
	ok := false
	if Backend == "rethinkdb" {
		res, ok = rethinkdb.GetFromDb(url)
	}
	return res, ok
}

func GetRoutes() []string {
	routes := []string{}
	if Backend == "rethinkdb" {
		routes = rethinkdb.GetRoutes()
	}
	return routes
}



func CommandsListener(comchan chan *datatypes.Command) {
	if conf.ListenToChangefeeds() == true {
		rethinkdb.CommandsListener(comchan)
	}
}

func PageChangesListener(c chan *datatypes.DataChanges) {
	if conf.ListenToChangefeeds() == true {
		rethinkdb.PageChangesListener(c)
	}
}
