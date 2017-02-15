package db

import (
	/*"fmt"
	"sync"
	r "gopkg.in/dancannon/gorethink.v2"
	postgresql "github.com/synw/microb/db/postgresql"*/
	"github.com/synw/microb/libmicrob/conf"
	"github.com/synw/microb/libmicrob/datatypes"
	"github.com/synw/microb/libmicrob/db/rethinkdb"
	//"github.com/synw/microb/utils"
	"github.com/synw/microb/libmicrob/metadata"
	
)

var Config = conf.GetConf()
var MainDb *datatypes.Database = metadata.GetMainDatabase()
var Backend = MainDb.Type


func GetFromUrl(url string)  (map[string]interface{}, bool) {
	res := make(map[string]interface{})
	ok := false
	if Backend == "rethinkdb" {
		res, ok = rethinkdb.GetFromDb(url)
	}
	return res, ok
}

func GetRoutes() []string {
	var routes []string
	if Backend == "rethinkdb" {
		routes = rethinkdb.GetRoutes()
	}
	return routes
}
