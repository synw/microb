package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/synw/centcom"
	"github.com/synw/microb/libmicrob/events"
	"github.com/synw/microb/libmicrob/types"
	"github.com/synw/microb/services"
	"github.com/synw/terr"
	"github.com/ventu-io/go-shortid"
	"time"
)

func isValid(command *types.Command) bool {
	is_valid := false
	for _, com := range services.ValidCommands {
		if com == command.Name {
			is_valid = true
			break
		}
	}
	return is_valid
}

func Run(payload interface{}, cli *centcom.Cli, server *types.Server) {
	cmd, exec := CmdFromPayload(payload)
	events.Cmd(cmd)
	if exec == false {
		return
	}
	if isValid(cmd) == false {
		fmt.Println("Invalid command", cmd)
		return
	}
	c := make(chan *types.Command)
	go services.Dispatch(cmd, c)
	select {
	case com := <-c:
		events.CmdExec(cmd)
		tr := sendCommand(com, cli, server)
		if tr != nil {
			msg := "Error executing the " + cmd.Name + " command"
			events.Err(cmd.Service, cmd.From, msg, tr.ToErr())
		}
		close(c)
	}
}

func CmdFromPayload(payload interface{}) (*types.Command, bool) {
	pl := payload.(map[string]interface{})
	status := pl["Status"].(string)
	name := pl["Name"].(string)
	serv := pl["Service"].(string)
	var s *types.Service
	found := false
	for i, srv := range services.All {
		if srv.Name == serv {
			s = services.All[i]
			found = true
			break
		}
	}
	cmd := &types.Command{}
	if !found {
		err := errors.New("Service " + serv + "not found")
		tr := terr.New("cmd.CmdFromPayload", err)
		cmd.Status = "error"
		cmd.Trace = tr
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

func sendCommand(command *types.Command, cli *centcom.Cli, server *types.Server) *terr.Trace {
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
	_, err = cli.Http.Publish(server.CmdChanOut, payload)
	if err != nil {
		trace := terr.New("commands.SendCommand", err)
		return trace
	}
	return nil
}

func New(name string, service string, from string, reason string, args ...interface{}) *types.Command {
	id, _ := shortid.Generate()
	date := time.Now()
	status := "pending"
	var tr *terr.Trace
	var rvs []interface{}
	var err_msg string
	command := &types.Command{
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
