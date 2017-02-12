package commands

import (
	"fmt"
	"sync"
	"os"
	"io/ioutil"
	"html/template"
	"github.com/acmacalister/skittles"
	"github.com/synw/microb/conf"
	"github.com/synw/microb/db"
	"github.com/synw/microb/db/datatypes"
	"github.com/synw/microb/utils"
)


var Config = conf.GetConf()

// commands
func reparseTemplates(c chan bool) {
	utils.PrintEvent("command", "Reparsing templates")
	template.Must(template.New("view.html").ParseFiles("templates/view.html", "templates/routes.js"))
	c <- true
}

func updateRoutes(c chan bool) {
	var routestab []string
	// hit db
	routestab = db.GetRoutes()
	var routestr string
	var route string
	for i := range(routestab) {
		route = routestab[i]
		routestr = routestr+fmt.Sprintf("page('%s', function(ctx, next) { loadPage('/x%s') } );", route, route)
	}
	utils.PrintEvent("command", "Rebuilding client side routes")
    str := []byte(routestr)
    err := ioutil.WriteFile("templates/routes.js", str, 0644)
    if err != nil {
        panic(err)
        c <- false
    } else {
    	c <- true
    }
 	return
}

// handlers
func isValid(command *datatypes.Command) bool {
	valid_commands := []string{"update_routes", "reparse_templates", "init_pg"}
	is_valid := false
	for _, com := range(valid_commands) {
		if (com == command.Name) {
			is_valid = true
			break
		}
	}
	return is_valid
}

func RunCommand(command *datatypes.Command, c chan bool, to_db bool) {
	if  isValid(command) == false {
		msg := "Unknown command: "+command.Name
		utils.PrintEvent("error", msg)
		c <- false
	}
	if command.From == "terminal" {
		msg := "Sending command "+skittles.BoldWhite(command.Name)+" to the server "
		msg = msg+Config["domain"].(string)+" from "+command.From
		if command.Reason != "nil" {
			msg = msg+". Reason: "+command.Reason
		}
		utils.PrintEvent("event", msg)
	}
	// run command
	if (command.Name == "update_routes") {
		go updateRoutes(c)
	} else if (command.Name == "reparse_templates") {
		go reparseTemplates(c)
	}
	// save in db
	if to_db == true {
		go db.SaveCommand(command)
	}
}

func RunCommandAndExit(command *datatypes.Command, wg *sync.WaitGroup, c chan bool) {
	defer os.Exit(0)
	RunCommand(command, c, false)
	go db.SaveCommandWait(command, wg)
	utils.PrintEvent("nil", "Exiting")
}
