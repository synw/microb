package mutate

import (
	"strconv"
	"errors"
    "github.com/synw/microb/libmicrob/state"
    "github.com/synw/microb/libmicrob/events"
    "github.com/synw/microb/libmicrob/http_handlers"
    "github.com/synw/microb/libmicrob/db"
)


func StartServer() error {
	if state.Server.Running == true {
		err := errors.New("Server is already started")
		events.State("mutate.StartServer", err.Error())
		return err
	}
	if state.Server.PagesDb.Name == "" {
		err := errors.New("Database pages is not configured: can't start server")
		events.State("mutate.StartServer", err.Error())
		return err
	}
	go state.Server.RunningServer.ListenAndServe()
	state.Server.Running = true
	if state.Verbosity > 0 {
		events.State("mutate.StartServer", "Server has started")
		http_handlers.StartMsg()
	}
	return nil	
}

func KillServer() error {
	if state.Server.Running == false {
		err := errors.New("Server is already stopped")
		events.State("mutate.StartServer", err.Error())
		return err
	}
	go state.Server.RunningServer.Close()
	state.Server.Running = false
	if state.Verbosity > 0 {
		events.State("mutate.StartServer", "Server has stopped")
	}
	return nil	
}

func Verbosity(lvl string) string {
	v, _ := strconv.Atoi(lvl)
	state.Verbosity = v
	msg := "Verbosity is set to "+lvl
	events.State("mutate.Verbosity", msg)
	return msg
}

func Debug(lvl string) (string, error) {
	var msg string
	if lvl == "on" {
		state.Debug = true
		msg = "Debug is on"
		events.State("mutate.Debug", msg)
	} else if lvl == "off" {
		state.Debug = false
		msg = "Debug is off"
		events.State("mutate.Debug", msg)
	} else {
		msg = "Invalid value for set debug: "+lvl
		events.ErrMsg("mutate.Debug", msg)
		err := errors.New(msg)
		return msg, err
	}
	return msg, nil
}

func PagesDb(name string) (string, error) {
	var msg string
	database, err := state.GetDbFromConf(name)
	if err != nil {
		return msg, err
	}
	err = db.SwitchDb("pages", database)
	if err != nil {
		return msg, err
	}
	state.Server.PagesDb = database
	msg = "Database pages is now set to "+name+" ("+database.Type+")"
	events.State("mutate.PagesDb", msg)
	return msg, nil
}
