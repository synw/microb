package db

import (
	"errors"
	"fmt"
	"github.com/synw/terr"
	"github.com/synw/microb/libmicrob/datatypes"
	"github.com/synw/microb/libmicrob/db/rethinkdb"
	"github.com/synw/microb/libmicrob/state"
	"github.com/synw/microb/libmicrob/conf"
	"github.com/synw/microb/libmicrob/events"
)


func GetByUrl(url string)  (*datatypes.Document, bool, *terr.Trace) {
	var r *datatypes.Document
	f := false
	if state.DocDb.Running == false {
		err := errors.New("Documents database is down")
		tr := terr.New("db.GetByUrl", err)
		return r, f, tr
	}
	if state.DocDb.Type == "rethinkdb" {
		res, found, tr := rethinkdb.GetByUrl(url)
		if tr != nil {
			return r, f, tr
		}
		return res, found, nil
	}
	err := errors.New("The database did not return any result")
	tr := terr.New("db.GetFromUrl", err)
	return r, f, tr
}

func InitDb(dev string) *terr.Trace {
	db, _ := conf.GetDefaultDb(dev)
	tr := CheckDb(db)
	if tr != nil {
		state.DocDb = nil
		events.Msg("state", "db.InitDb", "Disabling documents database")
		return tr
	}
	state.DocDb = db
	state.DocDb.Running = true
	if state.Verbosity > 0 {
		events.Msg("state", "db.InitDb", "Document database running")
	}
	if state.Verbosity > 1 {
		msg := "Database "+db.Name+" ( "+db.Type+") is up at "+db.Addr
		fmt.Println(terr.Ok(msg))
	}
	return nil
}

func CheckDb(db *datatypes.Database) *terr.Trace {
	err := errors.New("Unknown database type")
	var tr *terr.Trace
	if db.Type == "rethinkdb" {
		if db != nil {
			if state.Verbosity > 0 {
				fmt.Println("Connecting to database "+db.Name+" ...")
			}
			tr := rethinkdb.InitDb(db)
			if tr != nil {
				return tr
			}
		}
	} else {
		err = errors.New("Database type not implemented")
		tr = terr.New("db.rethinkdb.CheckDb", err)
		events.Err("error", "db.rethinkdb.CheckDb", tr)
	}
	return tr
}
