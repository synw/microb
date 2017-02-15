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
 	return nil
}

func ping() error {
	fmt.Println("PONG")
	return nil
}


// handlers
func RunCommand(command *datatypes.Command) error {
	events.New("command", "RunCommand()", command.Name)
	if  commands_methods.IsValid(command) == false {
		msg := "Unknown command: "+command.Name
		command.Status = "error"
		command.Error = errors.New(msg)
		return command.Error
	}
	// run command
	if (command.Name == "update_routes") {
		err := updateRoutes()
		if err != nil { 
			command.Status = "error" 
			command.Error = err
			return command.Error
		} else {
			command.Status = "success"
		}
	} else if (command.Name == "reparse_templates") {
		err := reparseTemplates()
		if err != nil { command.Status = "error" } else {
			command.Status = "success"
			command.Error = err
			return command.Error
		}
	} else if (command.Name == "ping") {
		err := ping()
		if err != nil { 
			command.Status = "error"
			command.Error = err
			return command.Error
		} else {
			command.Status = "success"
		}
	}
	/*else if (command.Name == "syncdb") {
		go db.ImportPagesFromMainDb(command.Values.(string))
	}*/
	// save in db
	//go db.SaveCommand(command)
	return nil
}

func Run(name string, from string, reason string) {
	command := &datatypes.Command{name, from, reason, time.Now(), "pending", nil}
	err := RunCommand(command)
	if err != nil {
		fmt.Println("Error executing command", name, ":", err)
	} else {
		fmt.Println("Command", name, "successfull")
	}
	return
}
