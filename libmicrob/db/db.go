package db

import (
	/*"fmt"
	"sync"
	r "gopkg.in/dancannon/gorethink.v2"
	postgresql "github.com/synw/microb/db/postgresql"*/
	"errors"
	"github.com/synw/microb/libmicrob/conf"
	"github.com/synw/microb/libmicrob/datatypes"
	"github.com/synw/microb/libmicrob/db/rethinkdb"
	"github.com/synw/microb/libmicrob/metadata"
	"github.com/synw/microb/libmicrob/events"
	
)

var Config = conf.GetConf()
var MainDb *datatypes.Database = metadata.GetMainDatabase()
var Backend = MainDb.Type


func ReportStatus() (map[string]interface{}, error) {
	var err error
	status := make(map[string]interface{})
	if Backend == "rethinkdb" {
		status, err = rethinkdb.ReportStatus()
	}
	if err != nil {
		return status, err
	}
	return status, nil
}

func GetFromUrl(url string)  (*datatypes.Page, bool, error) {
	var r *datatypes.Page
	f := false
	if Backend == "rethinkdb" {
		res, found, err := rethinkdb.GetFromUrl(url)
		if err != nil {
			events.Error("db.GetFromUrl", err)
			return r, f, nil
		}
		return res, found, nil
	}
	err := errors.New("The database did not return any result")
	events.Error("db.GetFromUrl", err)
	return r, f, nil
}

func GetRoutes() []string {
	var routes []string
	if Backend == "rethinkdb" {
		routes = rethinkdb.GetRoutes()
	}
	return routes
}
