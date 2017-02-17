package commands

import (
	"fmt"
	"io/ioutil"
	"html/template"
	"errors"
	"time"
	"github.com/synw/microb/libmicrob/conf"
	"github.com/synw/microb/libmicrob/db"
	"github.com/synw/microb/libmicrob/datatypes"
	"github.com/synw/microb/libmicrob/datatypes/encoding"
	"github.com/synw/microb/libmicrob/events"
	"github.com/synw/microb/libmicrob/commands/methods"
)


var Config = conf.GetConf()

// commands
func reparseTemplates() error {
	template.Must(template.New("view.html").ParseFiles("templates/view.html", "templates/routes.js"))
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
        panic(err)
        return err
    }
    // auto reparse templates if the routes change
    cmd := New("reparse_templates", "commands.updateRoutes()", "Routes update")
    go Run(cmd)
 	return nil
}

// utilities
func handleCommandError(command *datatypes.Command, err error, c chan *datatypes.Command) {
	if err != nil { 
		command.Status = "error" 
		command.Error = err
		c <- command
	}
	return
}

func GetCommandFromPayload(message *datatypes.WsIncomingMessage, broker string) *datatypes.Command {
	data := message.Data
	name := data["name"].(string)
	reason := ""
	from := broker
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

func HandleCommandFeedback(command *datatypes.Command) {
	if command.Status == "error" {
		events.Error("command_execution", command.Error)
	}
	events.PrintCommandReport(command)
}

// constructors
func New(name string, from string, reason string) *datatypes.Command {
	now := time.Now()
	var args[]interface{}
	var rv []string
	id := encoding.GenerateId()
	cmd := &datatypes.Command{id, name, from, reason, now, args, "pending", nil, rv}
	return cmd
}
func NewWithArgs(name string, from string, reason string, args []interface{}) *datatypes.Command {
	now := time.Now()
	var rv []string
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
			command.Status = "success"
		}
	} else if (command.Name == "reparse_templates") {
		err := reparseTemplates()
		if err != nil { 
			handleCommandError(command, err, c)
			return
		} else {
			command.Status = "success"
		}
	} else if (command.Name == "ping") {
		command.ReturnValues = append(command.ReturnValues, "PONG")
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
			if Config["verbosity"].(int) > 0 {
				HandleCommandFeedback(cmd)
			}
	}
}
