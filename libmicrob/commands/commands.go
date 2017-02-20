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
	"github.com/synw/microb/libmicrob/metadata"
	"github.com/synw/microb/libmicrob/commands/methods"
	"github.com/synw/microb/libmicrob/http_handlers"
	
)


// commands
func reparseTemplates() error {
	http_handlers.ReparseTemplates()
	return nil
}

func updateRoutes() error {
	// hit db
	routestab := db.GetRoutes()
	var routestr string
	var route string
	for i := range(routestab) {
		route = routestab[i]
		routestr = routestr+fmt.Sprintf("page('%s', function(ctx, next) { loadPage('/x%s') } );", route, route)
	}
    str := []byte(routestr)
    err := ioutil.WriteFile("templates/routes.js", str, 0644)
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

func HandleCommandFeedback(command *datatypes.Command) {
	if command.Status == "error" {
		events.Error("command_execution", command.Error)
	}
	if metadata.GetVerbosity() > 0 {
		events.PrintCommandReport(command)
	}
}

func GetCommandFromPayload(message *datatypes.WsIncomingMessage, broker string) *datatypes.Command {
	data := message.Data
	name := data["name"].(string)
	reason := ""
	from := broker
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
		err := updateRoutes()
		if err != nil {
			handleCommandError(command, err, c)
			return
		} else {
			go updateRoutes()
			command.Status = "success"
		}
	} else if (command.Name == "reparse_templates") {
		err := reparseTemplates()
		if err != nil { 
			handleCommandError(command, err, c)
			return
		} else {
			go reparseTemplates()
			command.Status = "success"
		}
	} else if (command.Name == "ping") {
		command.ReturnValues = append(command.ReturnValues, "PONG")
		command.Status = "success"
	} else if (command.Name == "db_status") {
		status, err := db.ReportStatus()
		if err != nil { 
			handleCommandError(command, err, c)
			return
		}
		line := "------------------------------"
		version := "\n "+line+"\n - Version: "+status["version"].(string)+"\n"
		cache_size_mb := "- Cache size: "+strconv.FormatFloat(status["cache_size_mb"].(float64), 'f', 2, 64)+" Mb\n"
		time_started := line+"\nStarted since "+metadata.FormatTime(status["time_started"].(time.Time))
		command.ReturnValues = append(command.ReturnValues, version)
		command.ReturnValues = append(command.ReturnValues, cache_size_mb)
		command.ReturnValues = append(command.ReturnValues, time_started)
		command.Status = "success"
	} else if (command.Name == "routes") {
		routes := db.GetRoutes()
		var rvs []interface{}
		for _, route := range(routes) {
			rvs = append(rvs, route)
		} 
		command.ReturnValues = rvs
		command.Status = "success"
	}
	c <- command
	/*else if (command.Name == "syncdb") {
		go db.ImportPagesFromMainDb(command.Values.(string))
	}*/
	// save in db
	//go db.SaveCommand(command)
	return
}

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
