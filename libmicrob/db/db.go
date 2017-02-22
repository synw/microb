package db

import (
	/*"fmt"
	"sync"
	postgresql "github.com/synw/microb/db/postgresql"*/
	"errors"
	"github.com/synw/microb/libmicrob/datatypes"
	"github.com/synw/microb/libmicrob/db/rethinkdb"
	"github.com/synw/microb/libmicrob/events"
	"github.com/synw/microb/libmicrob/state"
	
)


var db_type string
var misconf = errors.New("Database for pages is not configurated")

func init() {
	setPagesDb()
}

func ReportStatus() (map[string]interface{}, error) {
	var err error
	status := make(map[string]interface{})
	if db_type == "" {
		return status, misconf
	}
	if db_type == "rethinkdb" {
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
	if db_type == "" {
		return r, f, misconf
	}
	if db_type == "rethinkdb" {
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

func GetRoutes() ([]string, error) {
	var routes []string
	if db_type == "" {
		return routes, misconf
	}
	if db_type == "rethinkdb" {
		routes = rethinkdb.GetRoutes()
	}
	return routes, nil
}

func setPagesDb() {
	if state.Server.PagesDb != nil {
		db_type = state.Server.PagesDb.Type
	}
}

/*
// constructor
func New(dbtype, name, host, port, user, pwd) {
	database := &datatypes.Database{dbtype, name, host, port, user, pwd}
	return database
}*/
