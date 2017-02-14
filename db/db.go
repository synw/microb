package db

import (
	/*"fmt"
	"sync"
	r "gopkg.in/dancannon/gorethink.v2"
	postgresql "github.com/synw/microb/db/postgresql"*/
	"github.com/synw/microb/conf"
	"github.com/synw/microb/datatypes"
	"github.com/synw/microb/db/rethinkdb"
	//"github.com/synw/microb/utils"
	
)

var Config = conf.GetConf()
var MainDb *datatypes.Database = conf.GetMainDatabase()
var Backend = MainDb.Type


func GetFromUrl(url string)  (map[string]interface{}, bool) {
	res := make(map[string]interface{})
	ok := false
	if Backend == "rethinkdb" {
		res, ok = rethinkdb.GetFromDb(url)
	}
	return res, ok
}