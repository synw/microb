package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	color "github.com/acmacalister/skittles"
	"github.com/synw/microb/libmicrob/datatypes"
	"github.com/synw/microb/libmicrob/events"
	"github.com/synw/microb/libmicrob/state"
	"github.com/synw/microb/services"
	"github.com/synw/terr"
	"github.com/ventu-io/go-shortid"
	"time"
)

func isValid(command *datatypes.Command) bool {
	is_valid := false
	for _, com := range services.ValidCommands {
		if com == command.Name {
			is_valid = true
			break
		}
	}
	return is_valid
}

func Run(payload interface{}) {
	cmd, exec := CmdFromPayload(payload)
	if exec == false {
		return
	}
	if isValid(cmd) == false {
		fmt.Println("Invalid command", cmd)

		return
	}
	c := make(chan *datatypes.Command)
	go services.Dispatch(cmd, c)
	select {
	case com := <-c:
		status := com.Status
		if status == "error" {
			status = color.BoldRed("error")
			if state.Verbosity > 0 {
				fmt.Println(" ->", status, com.Trace.Format())
			}
		} else if status == "success" {
			status = color.Green("success")
			if state.Verbosity > 0 {
				fmt.Println(" ->", status, com.ReturnValues)
			}
		}
		tr := sendCommand(com)
		if tr != nil {
			events.Err("error", cmd.Name, tr)
		}
		close(c)
	}
}

func CmdFromPayload(payload interface{}) (*datatypes.Command, bool) {
	pl := payload.(map[string]interface{})
	status := pl["Status"].(string)
	name := pl["Name"].(string)
	serv := pl["Service"].(string)
	var s *datatypes.Service
	found := false
	for i, srv := range services.All {
		if srv.Name == serv {
			s = services.All[i]
			found = true
			break
		}
	}
	cmd := &datatypes.Command{}
	if !found {
		terr.Debug("Service not found", s)

		return cmd, false
	}
	cmd.Service = s.Name
	from := pl["From"].(string)
	reason := pl["Reason"].(string)
	var args []interface{}
	if pl["Args"] != nil {
		args = pl["Args"].([]interface{})
	}
	if args != nil {
		cmd = New(name, s.Name, from, reason, args)
	} else {
		cmd = New(name, s.Name, from, reason)
	}
	if pl["ErrMsg"] != "" {
		msg := pl["ErrMsg"].(string)
		err := errors.New(msg)
		cmd.Trace = terr.New("cmd.CmdFromPayload", err)
		cmd.ErrMsg = msg
	}
	if pl["ReturnValues"] != nil {
		cmd.ReturnValues = pl["ReturnValues"].([]interface{})
	}
	cmd.Status = status
	if status != "pending" {
		return cmd, false
	}
	return cmd, true
}

func printJson(b []byte) ([]byte, error) {
	var out bytes.Buffer
	err := json.Indent(&out, b, "", "  ")
	return out.Bytes(), err
}

func sendCommand(command *datatypes.Command) *terr.Trace {
	if command.Trace != nil {
		command.ErrMsg = command.Trace.Formatc()
		command.Status = "error"
	} else {
		command.Status = "success"
	}
	payload, err := json.Marshal(command)
	if err != nil {
		msg := "Unable to marshall json: " + err.Error()
		err := errors.New(msg)
		trace := terr.New("commands.SendCommand", err)
		return trace
	}
	_, err = state.Cli.Http.Publish(state.Server.CmdChanOut, payload)
	if err != nil {
		trace := terr.New("commands.SendCommand", err)
		return trace
	}
	return nil
}

func New(name string, service string, from string, reason string, args ...interface{}) *datatypes.Command {
	id, _ := shortid.Generate()
	date := time.Now()
	status := "pending"
	var tr *terr.Trace
	var rvs []interface{}
	var err_msg string
	command := &datatypes.Command{
		id,
		service,
		name,
		from,
		reason,
		date,
		args,
		status,
		tr,
		err_msg,
		rvs,
	}
	return command
}
