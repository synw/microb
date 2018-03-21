package cmds

import (
	"encoding/json"
	"errors"
	"github.com/SKAhack/go-shortid"
	"github.com/synw/microb/libmicrob/events"
	"github.com/synw/microb/libmicrob/msgs"
	"github.com/synw/microb/libmicrob/types"
	"github.com/synw/terr"
)

var g = shortid.Generator()

func checkServiceCmd(cmd *types.Cmd, state *types.State) (*types.Cmd, bool) {
	isValid := false
	for _, srv := range state.Services {
		if srv.Name == cmd.Name {
			cmd.Service = cmd.Name
			cmd.Name = cmd.Args[0].(string)
			if len(cmd.Args) > 1 {
				cmd.Args = cmd.Args[1:]
			}
			isValid = true
			break
		}
	}
	return cmd, isValid
}

func getCmd(cmd *types.Cmd, state *types.State) (*types.Cmd, bool) {
	for sname, srv := range state.Services {
		if sname == cmd.Service {
			for cname, scmd := range srv.Cmds {
				if cname == cmd.Name {
					cmd.Exec = scmd.Exec
					return cmd, true
				}
			}
		}
	}
	return cmd, false
}

func Run(payload interface{}, state *types.State) {
	cmd := ConvertPayload(payload)
	msgs.Debug("111", cmd.Args, cmd.Name, cmd.Service)

	cmd, isValid := getCmd(cmd, state)
	if isValid == false {
		msgs.Error("Invalid command " + cmd.Name)
		return
	}
	msgs.Debug("222", cmd.Args, cmd.Service)

	exec := state.Cmds[cmd.Name].Exec.(func(*types.Cmd, chan *types.Cmd, ...interface{}))
	events.Cmd(cmd)

	c := make(chan *types.Cmd)
	go exec(cmd, c, state)
	select {
	case com := <-c:
		events.CmdExec(cmd)
		// set to interface to be able to marshall json
		com.Exec = nil
		tr := sendCommand(com, state)
		if tr != nil {
			msg := "Error executing the " + cmd.Name + " command"
			//events.Err(cmd.Service, cmd.From, msg, tr.ToErr())
			tr.Print()
			msgs.Error(msg)
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
		Id:      g.Generate(),
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

func sendCommand(cmd *types.Cmd, state *types.State) *terr.Trace {
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
	_, err = state.Cli.Http.Publish(state.WsServer.CmdChanOut, payload)
	if err != nil {
		trace := terr.New("commands.SendCommand", err)
		return trace
	}
	return nil
}
