package commands

import (
	"fmt"
	"io/ioutil"
	"html/template"
	"errors"
	"github.com/synw/microb/libmicrob/conf"
	"github.com/synw/microb/libmicrob/db"
	"github.com/synw/microb/libmicrob/datatypes"
	"github.com/synw/microb/libmicrob/events"
	"github.com/synw/microb/libmicrob/commands/methods"
)


var Config = conf.GetConf()

func PrintFeedback(command *datatypes.Command) {
	if command.Status == "error" {
		events.Error("command_execution", command.Error)
	} else if command.Status == "success" {
		events.New("runtime_info", "command_feedback", "Command successfull")
	} else {
		msg := "Command status: "+command.Status
		events.New("runtime_info", "command_feedback", msg)
	}
}

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
 	return nil
}

func ping() error {
	fmt.Println("PONG")
	return nil
}

func handleCommandError(command *datatypes.Command, err error, c chan *datatypes.Command) {
	if err != nil { 
		command.Status = "error" 
		command.Error = err
		c <- command
	}
	return
}

// handlers
func Run(command *datatypes.Command, c chan *datatypes.Command) {
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
		err := ping()
		if err != nil { 
			handleCommandError(command, err, c)
			return
		} else {
			command.Status = "success"
		}
	}
	c <- command
	/*else if (command.Name == "syncdb") {
		go db.ImportPagesFromMainDb(command.Values.(string))
	}*/
	// save in db
	//go db.SaveCommand(command)
	return
}
