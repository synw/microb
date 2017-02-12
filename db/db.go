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

func commandsTransports() []string {
	var ts []string
	cts := Config["commands_transport"].([]string)
	// check for defaults
	is_default := false
	for _, transp := range cts {
		if transp == "default" {
			is_default = true
			if Config["db_type"].(string) == "rethinkdb" {
				ts = []string{"changefeeds"}
			}
		}
	}
	if is_default == false {
		ts = cts
	}
	/*var ts_str string
	for _, transport := range ts {
		ts_str = ts_str+" "+transport
	} 
	utils.PrintEvent("info", "Transports used for commands:"+ts_str)*/
	return ts
}

func listenToChangefeeds() bool {
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

func CommandsListener(comchan chan *datatypes.Command) {
	if listenToChangefeeds() == true {
		rethinkdb.CommandsListener(comchan)
	}
}

func PageChangesListener(c chan *datatypes.DataChanges) {
	if listenToChangefeeds() == true {
		rethinkdb.PageChangesListener(c)
	}
}
