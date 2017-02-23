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


var pages_db string
var misconf = errors.New("Database for pages is not configurated")

func InitDb() error {
	if state.Server.PagesDb.Name != "" {
		pages_db = state.Server.PagesDb.Type
	} else {
		err := errors.New("No database set for pages")
		return err
	}
	err := CheckDb(state.Server.PagesDb)
	if err != nil {
		state.Server.PagesDb = nil
		state.Server.PagesDb.Running = false
		events.State("db.InitDb", "Disabling pages database")
		return err
	}
	state.Server.PagesDb.Running = true
	if state.Debug == true {
		events.Debug("Pages db is now running")
	}
	return nil
}

func CheckDb(database *datatypes.Database) error {
	err := errors.New("Unknown database type")
	if database.Type == "rethinkdb" {
		if database != nil {
			err = rethinkdb.InitDb(state.Server.PagesDb)
			if err != nil {
				events.ErrMsg("db.InitDb", "Database connection error")
				return nil
			}
		}
	} else {
		err = errors.New("Not implemented")
		events.Error("db.rethinkdb.CheckDb", err)
	}
	return err
}

func SwitchDb(role string, database *datatypes.Database) error {
	var err error
	if role == "pages" { 
		if database.Type == "rethinkdb" {
			err = rethinkdb.SwitchDb(role, database)
			if err != nil {
				events.Error("db.SwitchDb", err)
				return err
			} else {
				database.Running = true
				events.New("info", "db.SwitchDb", "Database pages is running")
			}
		} else {
			err = errors.New("Not implemented")
			events.Error("db.rethinkdb.CheckDb", err)
		}
	} else {
		err = errors.New("Not implemented")
		events.Error("db.rethinkdb.SwitchDb", err)
	}
	return nil
}

func ReportStatus() (map[string]interface{}, error) {
	var err error
	status := make(map[string]interface{})
	if pages_db == "" {
		return status, misconf
	}
	if pages_db == "rethinkdb" {
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
	if pages_db == "" {
		return r, f, misconf
	}
	if pages_db == "rethinkdb" {
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
	if pages_db == "" {
		return routes, misconf
	}
	if pages_db == "rethinkdb" {
		routes = rethinkdb.GetRoutes()
	}
	return routes, nil
}

/*
// constructor
func New(dbtype, name, host, port, user, pwd) {
	database := &datatypes.Database{dbtype, name, host, port, user, pwd}
	return database
}*/
