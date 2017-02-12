package commands

import (
	"sync"
	"os"
	"github.com/acmacalister/skittles"
	"github.com/synw/microb/conf"
	"github.com/synw/microb/db"
	"github.com/synw/microb/db/datatypes"
	"github.com/synw/microb/utils"
)


var Config = conf.GetConf()

var transport = Config["commands_transport"].([]string)

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

func RunCommand(command *datatypes.Command, wg *sync.WaitGroup) bool {
	if  isValid(command) == false {
		msg := "Unknown command: "+command.Name
		utils.PrintEvent("error", msg)
		wg.Done()
		return false
	}
	if command.From == "terminal" {
		msg := "Sending command "+skittles.BoldWhite(command.Name)+" to the server "
		msg = msg+Config["domain"].(string)+" from "+command.From
		if command.Reason != "nil" {
			msg = msg+". Reason: "+command.Reason
		}
		utils.PrintEvent("event", msg)
	}
	// save in db
	go db.SaveCommand(command, wg)
	return true
}

func RunCommandAndExit(command *datatypes.Command, wg *sync.WaitGroup) {
	RunCommand(command, wg)
	utils.PrintEvent("nil", "Exiting")
	os.Exit(0)
}
