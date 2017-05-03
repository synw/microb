package db

import (
	"errors"
	"github.com/synw/microb/libmicrob/datatypes"
	//"github.com/synw/microb/libmicrob/db/rethinkdb"
	"github.com/synw/microb/libmicrob/state"
	"github.com/synw/microb/libmicrob/conf"
	"github.com/synw/microb/libmicrob/events"
	
)

func InitDb(dev string) error {
	db, _ := conf.GetDefaultDb(dev)
	/*err := CheckDb(db)
	if err != nil {
		state.DocDb = nil
		state.DocDb.Running = false
		events.State("db.InitDb", "Disabling documents database")
		return err
	}*/
	state.DocDb = db
	state.DocDb.Running = true
	if state.Debug == true {
		events.New("Document db is now running", "db:InitDb")
	}
	return nil
}

func CheckDb(database *datatypes.Database) error {
	err := errors.New("Unknown database type")
	/*if database.Type == "rethinkdb" {
		if database != nil {
			err = rethinkdb.InitDb(state.DocDb)
			if err != nil {
				events.ErrMsg("db.InitDb", "Database connection error")
				return nil
			}
		}
	} else {
		err = errors.New("Not implemented")
		events.Error("db.rethinkdb.CheckDb", err)
	}*/
	return err
}
