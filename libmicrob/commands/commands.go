package commands

import (
	"fmt"
	"io/ioutil"
	"errors"
	"time"
	"strconv"
	//"github.com/shirou/gopsutil/mem"
	"github.com/synw/microb/libmicrob/db"
	"github.com/synw/microb/libmicrob/datatypes"
	"github.com/synw/microb/libmicrob/datatypes/encoding"
	"github.com/synw/microb/libmicrob/events"
	"github.com/synw/microb/libmicrob/events/format"
	"github.com/synw/microb/libmicrob/state"
	"github.com/synw/microb/libmicrob/state/mutate"
	"github.com/synw/microb/libmicrob/commands/methods"
	"github.com/synw/microb/libmicrob/http_handlers"
	
)


func Run(command *datatypes.Command) {
	c := make(chan *datatypes.Command)
	go runCommand(command, c)
	select {
		case cmd := <- c:
			close(c)
			// process command results
			HandleCommandFeedback(cmd)
	}
}

func HandleCommandFeedback(command *datatypes.Command) {
	if command.Status == "error" {
		events.Error("command_execution", command.Error)
	}
	if state.Verbosity > 0 {
		events.PrintCommandReport(command)
	}
}

func GetCommandFromPayload(message *datatypes.WsIncomingMessage, from string) *datatypes.Command {
	data := message.Data
	name := data["name"].(string)
	reason := ""
	if data["from"] != nil {
		from = data["from"].(string)
	}
	var command *datatypes.Command
	if data["reason"].(string) != "" {
		reason = data["reason"].(string)
	}
	if data["args"] != nil {
		args := data["args"].([]interface{})
		command = NewWithArgs(name, from, reason, args)
	} else {
		command = New(name, from, reason)
	}
	return command
}

// constructors
func New(name string, from string, reason string) *datatypes.Command {
	now := time.Now()
	var args[]interface{}
	var rv []interface{}
	id := encoding.GenerateId()
	cmd := &datatypes.Command{id, name, from, reason, now, args, "pending", nil, rv}
	return cmd
}

func NewWithArgs(name string, from string, reason string, args []interface{}) *datatypes.Command {
	now := time.Now()
	var rv []interface{}
	id := encoding.GenerateId()
	cmd := &datatypes.Command{id, name, from, reason, now, args, "pending", nil, rv}
	return cmd
}

// commands
func reparseTemplates() error {
	http_handlers.ReparseTemplates()
	return nil
}

func updateRoutes() error {
	// hit db
	routestab, err := db.GetRoutes()
	if err != nil {
		events.Error("commands.updateRoutes", err)
	}
	var routestr string
	var route string
	for i := range(routestab) {
		route = routestab[i]
		routestr = routestr+fmt.Sprintf("page('%s', function(ctx, next) { loadPage('/x%s') } );", route, route)
	}
    str := []byte(routestr)
    err = ioutil.WriteFile("templates/routes.js", str, 0644)
    if err != nil {
        events.Error("commands.updateRoutes", err)
        return err
    }
    // auto reparse templates if the routes change
   	cmd := New("reparse_templates", "commands.updateRoutes()", "Routes update")
    go Run(cmd)
 	return nil
}

func dbStatus() (map[string]interface{}, error) {
	status, err := db.ReportStatus()
	if err != nil {
		return status, err
	}
	return status, nil
}
/*
func sysInfo(*datatypes.Command) *datatypes.Command {
	v, _ := mem.VirtualMemory()
	
}*/

// utilities
func handleCommandError(command *datatypes.Command, err error, c chan *datatypes.Command) {
	if err != nil { 
		command.Status = "error" 
		command.Error = err
		c <- command
	}
	return
}
/*
func NewWithData(name string, from string, reason string, data []interface{}) *datatypes.Command {
	now := time.Now()
	cmd := &datatypes.Command{name, from, reason, now, "pending", nil, data}
	return cmd
}
*/
// handlers
func runCommand(command *datatypes.Command, c chan *datatypes.Command) {
	events.New("command", "RunCommand()", command.Name)
	if  commands_methods.IsValid(command) == false {
		msg := "Unknown command: "+command.Name
		err := errors.New(msg)
		handleCommandError(command, err, c)
		return
	}
	if (command.Name == "update_routes") {
		// UPDATE ROUTES
		err := updateRoutes()
		if err != nil {
			handleCommandError(command, err, c)
			return
		} else {
			go updateRoutes()
			command.Status = "success"
		}
	} else if (command.Name == "reparse_templates") {
		// REPARSE TEMPLATES
		err := reparseTemplates()
		if err != nil { 
			handleCommandError(command, err, c)
			return
		} else {
			go reparseTemplates()
			command.Status = "success"
		}
	} else if (command.Name == "ping") {
		// PING
		command.ReturnValues = append(command.ReturnValues, "PONG")
		command.Status = "success"
	} else if (command.Name == "db_status") {
		status, err := db.ReportStatus()
		if err != nil { 
			handleCommandError(command, err, c)
			return
		}
		line := "------------------------------"
		version := "\n "+line+"\n - Version: "+status["version"].(string)
		cache_size_mb := "- Cache size: "+strconv.FormatFloat(status["cache_size_mb"].(float64), 'f', 2, 64)+" Mb"
		time_started := line+"\nStarted since "+format.FormatTime(status["time_started"].(time.Time))
		command.ReturnValues = append(command.ReturnValues, version)
		command.ReturnValues = append(command.ReturnValues, cache_size_mb)
		command.ReturnValues = append(command.ReturnValues, time_started)
		command.Status = "success"
	} else if (command.Name == "routes") {
		// ROUTES
		routes, err := db.GetRoutes()
		if err != nil {
			handleCommandError(command, err, c)
			return
		}
		var rvs []interface{}
		for _, route := range(routes) {
			rvs = append(rvs, route)
		}
		command.ReturnValues = rvs
		command.Status = "success"
	} else if command.Name == "state" {
		// STATE
		num_args := len(command.Args)
		if  num_args == 0 {
			msg := state.FormatState()
			var rvs []interface{}
			rvs = append(rvs, msg)
			command.ReturnValues = rvs
			command.Status = "success"
		} else if num_args > 1 {
			var rvs []interface{}
			msg := "Too many arguments to command state"
			rvs = append(rvs, msg)
			command.ReturnValues = rvs
			command.Status = "error"
		} else {
			if command.Args[0].(string) == "db" {
				// STATE DB
				var rvs []interface{}
				msg := "Database state:\n"
				msg = msg+" - Pages database is "+state.Server.PagesDb.Name+" ("+state.Server.PagesDb.Type+") at "+state.Server.PagesDb.Host+"\n"
				msg = msg+" - Hits database is "+state.Server.HitsDb.Name+" ("+state.Server.HitsDb.Type+") at "+state.Server.HitsDb.Host+"\n"
				msg = msg+" - Commands database is "+state.Server.CommandsDb.Name+" ("+state.Server.CommandsDb.Type+") at "+state.Server.CommandsDb.Host
				rvs = append(rvs, msg)
				command.ReturnValues = rvs
				command.Status = "success"
			}
		}
	} else if command.Name == "set" {
		// SET
		num_args := len(command.Args)
		if  num_args == 0 {
			msg := "Command set incoming from "+command.From+" with no arguments"
			events.ErrMsg("commands.run", msg)
		} else if num_args == 1 {
			msg := "Command set incoming from "+command.From+" with only one argument: "+command.Args[0].(string)
			events.ErrMsg("commands.run", msg)
		} else if num_args > 2 {
			var arguments string
			for _, a := range(command.Args) {
				arguments = arguments+" "+a.(string)
			}
			msg := "Command set incoming from "+command.From+" with more than two arguments:"+arguments
			events.ErrMsg("commands.run", msg)
		} else {
			cmd := command.Args[0]
			arg := command.Args[1]
			if cmd == "verbosity" {
				msg := mutate.Verbosity(arg.(string))
				var rvs []interface{}
				rvs = append(rvs, msg)
				command.ReturnValues = rvs
				command.Status = "success"
			} else if cmd == "debug" {
				msg, err := mutate.Debug(arg.(string))
				if err != nil {
					command.Status = "error"
					command.Error = err
				} else {
					var rvs []interface{}
					rvs = append(rvs, msg)
					command.Status = "success"
				}
			} 
		}
	}
	/*else if (command.Name == "db_list") {
		for role, db := range metadata.GetDatabasesAndRoles() {
			msg := database.Name
			command.ReturnValues = append(command.ReturnValues, )
		}
		command.Status = "success"
	}*/
	c <- command
	/*else if (command.Name == "syncdb") {
		go db.ImportPagesFromMainDb(command.Values.(string))
	}*/
	// save in db
	//go db.SaveCommand(command)
	return
}
