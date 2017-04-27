package cmd

import (
	"fmt"
	"time"
	"errors"
	"encoding/json"
	"github.com/ventu-io/go-shortid"
	"github.com/synw/terr"
	"github.com/synw/microb/libmicrob/datatypes"
	"github.com/synw/microb/libmicrob/state"
	"github.com/synw/microb/libmicrob/cmd/base"
)


func IsValid(command *datatypes.Command) bool {
	valid_commands := []string{"ping"}
	is_valid := false
	for _, com := range(valid_commands) {
		if (com == command.Name) {
			is_valid = true
			break
		}
	}
	return is_valid
}

func Run(payload interface{}) {
	cmd, exec := CmdFromPayload(payload)
	if (exec == false) {
		return
	}
	c := make(chan *datatypes.Command)
	go runCommand(cmd, c)
	select {
		case com := <- c:
			fmt.Println("RES", com)
			err := sendCommand(com)
			if err != nil {
				// todo
			}
			close(c)
	}
}

func CmdFromPayload(payload interface{}) (*datatypes.Command, bool) {
	pl := payload.(map[string]interface{})
	status := pl["Status"].(string)
	name := pl["Name"].(string)
	from := pl["From"].(string)
	reason := pl["Reason"].(string)
	cmd := &datatypes.Command{}
	var args []interface{}
	if pl["Args"] != nil {
		args = pl["Args"].([]interface{})
	}
	if args != nil {
		cmd = New(name, from, reason, args)
	} else {
		cmd = New(name, from, reason)
	}
	if pl["Error"] != nil {
		cmd.Error = pl["Error"].(error)
	}
	if pl["ReturnValues"] != nil {
		cmd.ReturnValues = pl["ReturnValues"].([]interface{})
	}
	if (status != "pending") {
		return cmd, false
	}
	return cmd, true
}

func sendCommand(command *datatypes.Command) *terr.Trace {
	payload, err := json.Marshal(command)
	if err != nil {
		msg := "Unable to marshall json: "+err.Error()
		err := errors.New(msg)
		trace := terr.New("commands.SendCommand", err)
		return trace
	}
	_, err = state.Cli.Http.Publish(state.Server.CmdChannel, payload)
	if err != nil {
		trace := terr.New("commands.SendCommand", err)
		return trace
	}
	return nil
}

func New(name string, from string, reason string, args ...interface{}) *datatypes.Command {
	id, _ := shortid.Generate()
	date := time.Now()
	status := "pending"
	var err error
	var rvs []interface{}
	command := &datatypes.Command{
		id,
		name,
		from,
		reason,
		date,
		args,
		status,
		err,
		rvs,
	}
	return command
}

func runCommand(cmd *datatypes.Command, c chan *datatypes.Command) {
	com := &datatypes.Command{}
	if cmd.Name == "ping" {
		com = base.Ping(cmd)
	}
	c <- com
	return
}
