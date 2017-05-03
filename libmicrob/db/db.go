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
	if state.Debug == true {
		events.Msg("state", "db.InitDb", "Document database is running")
	}
	if state.Verbosity > 1 {
		msg := "Database "+db.Name+" is up at "+db.Addr
		fmt.Println(terr.Ok(msg))
	}
	return nil
}

func CheckDb(db *datatypes.Database) *terr.Trace {
	err := errors.New("Unknown database type")
	var tr *terr.Trace
	if db.Type == "rethinkdb" {
		if db != nil {
			tr := rethinkdb.InitDb(db)
			if tr != nil {
				return tr
			}
		}
	} else {
		err = errors.New("Database type not implemented")
		tr = terr.New("db.rethinkdb.CheckDb", err)
		events.Error(tr)
	}
	return tr
}
