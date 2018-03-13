package cmd

import (
	//sid "github.com/SKAhack/go-shortid"
	"encoding/json"
	"errors"
	m "github.com/synw/microb/libmicrob"
	"github.com/synw/microb/libmicrob/events"
	"github.com/synw/microb/libmicrob/types"
	"github.com/synw/microb/services"
	"github.com/synw/terr"
	//"time"
)

func Run(payload interface{}) {
	cmd := ConvertPayload(payload)
	events.Cmd(cmd)
	/*if isValid(cmd) == false {
		fmt.Println("Invalid command", cmd)
		return
	}*/
	c := make(chan *types.Cmd)
	go services.Dispatch(cmd, c)
	select {
	case com := <-c:
		events.CmdExec(cmd)
		tr := sendCommand(com)
		if tr != nil {
			msg := "Error executing the " + cmd.Name + " command"
			//events.Err(cmd.Service, cmd.From, msg, tr.ToErr())
			m.Error(msg)
		}
		close(c)
	}
}

func ConvertPayload(payload interface{}) *types.Cmd {
	pl := payload.(map[string]interface{})
	status := pl["Status"].(string)
	name := pl["Name"].(string)
	serv := pl["Service"].(string)
	from := pl["From"].(string)
	var args []interface{}
	if pl["Args"] != nil {
		args = pl["Args"].([]interface{})
	}
	cmd := &types.Cmd{
		Name:    name,
		From:    from,
		Args:    args,
		Status:  status,
		Service: serv,
	}
	if args != nil {
		cmd.Args = args
	}
	if pl["ErrMsg"].(string) != "" {
		msg := pl["ErrMsg"].(string)
		err := errors.New(msg)
		cmd.Trace = terr.New("cmd.CmdFromPayload", err)
	}
	if pl["ReturnValues"] != nil {
		cmd.ReturnValues = pl["ReturnValues"].([]interface{})
	}
	cmd.Status = status
	return cmd
}

func sendCommand(cmd *types.Cmd) *terr.Trace {
	if cmd.Trace != nil {
		cmd.ErrMsg = cmd.Trace.Formatc()
		cmd.Status = "error"
	} else {
		cmd.Status = "success"
	}
	payload, err := json.Marshal(cmd)
	if err != nil {
		msg := "Unable to marshall json: " + err.Error()
		err := errors.New(msg)
		trace := terr.New("commands.SendCommand", err)
		return trace
	}
	_, err = m.Cli.Http.Publish(m.Server.CmdChanOut, payload)
	if err != nil {
		trace := terr.New("commands.SendCommand", err)
		return trace
	}
	return nil
}

/*
func isValid(cmd *types.Command) bool {
	is_valid := false
	for _, com := range ValidCommands {
		if com == cmd.Name {
			is_valid = true
			break
		}
	}
	return is_valid
}

func getAllValidCommands() []string {
	vc := []string{}
	for _, s := range All {
		vc = append(vc, All[s.Name].Cmds...)
	}
	return vc
}

func CmdIsValid(cmd *types.Command) bool {
	name := cmd.Name
	for _, c := range getAllValidCommands() {
		if c == name {
			return true
		}
	}
	return false
}
*/
