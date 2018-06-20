package cmds

import (
	"encoding/json"
	"errors"
	"github.com/SKAhack/go-shortid"
	//"github.com/davecgh/go-spew/spew"
	"github.com/synw/microb/events"
	"github.com/synw/microb/msgs"
	"github.com/synw/microb/types"
	"github.com/synw/terr"
	"time"
)

var g = shortid.Generator()

/*func Error(service string, msg string, cmd *types.Cmd, tr *terr.Trace) (cmd *types.Cmd) {
	events.Error(service, msg, tr)
	cmd.Status = "error"
	cmd.ErrMsg = tr.Format()
	cmd.Trace = tr
	cmd.Service = service
	return cmd
}*/

func Run(payload interface{}, state *types.State) {
	/*
		Runs a command from a json payload
	*/
	cmd := ConvertPayload(payload)
	cmd, isValid := getCmd(cmd, state)
	if isValid == false {
		msgs.Error("Invalid command " + cmd.Name)
		return
	}
	events.CmdIn(cmd)
	// execute the command
	c := make(chan *types.Cmd)
	cname := cmd.Service + "_" + cmd.Name
	if cmd.ExecCli != nil {
		exec := state.Cmds[cname].ExecCli.(func(*types.Cmd, chan *types.Cmd, ...interface{}))
		go exec(cmd, c, state)
	} else {
		exec := state.Cmds[cname].Exec.(func(*types.Cmd, chan *types.Cmd, ...interface{}))
		go exec(cmd, c, state)
	}
	select {
	case com := <-c:
		// check for errors
		if cmd.Trace != nil {
			msg := "Error executing the " + cmd.Name + " command "
			cmd.Trace.Add(msg)
			msg = msg + " from the " + cmd.Service + " service"
		}
		// send back the command
		events.CmdOut(cmd)
		// set the Exec to nil to be able to marshall json
		com.Exec = nil
		tr := sendCommand(com, state)
		if tr != nil {
			msg := "Error sending back the " + cmd.Name + " command"
			tr.Add(msg)
			events.Error("microb", msg, tr)
		}
		close(c)
	}
}

func ConvertPayload(payload interface{}) *types.Cmd {
	/*
		Converts a json payload to a command
	*/
	pl := payload.(map[string]interface{})
	status := pl["Status"].(string)
	name := pl["Name"].(string)
	serv := pl["Service"].(string)
	from := pl["From"].(string)
	dateStr := pl["Date"].(string)
	domain := pl["Domain"].(string)
	noLog := false
	if _, ok := pl["NoLog"]; ok {
		noLog = pl["NoLog"].(bool)
	}
	id := pl["Id"].(string)
	if id == "" {
		id = g.Generate()
	}
	date, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		tr := terr.New(err)
		events.Error("microb", "Can not parse date from json payload", tr)
	}
	var args []interface{}
	if pl["Args"] != nil {
		args = pl["Args"].([]interface{})
	}
	cmd := &types.Cmd{
		Id:      id,
		Name:    name,
		Date:    date,
		From:    from,
		Domain:  domain,
		Args:    args,
		Status:  status,
		Service: serv,
		NoLog:   noLog,
	}
	if args != nil {
		cmd.Args = args
	}
	if pl["ReturnValues"] != nil {
		cmd.ReturnValues = pl["ReturnValues"].([]interface{})
	}
	cmd.Status = status
	return cmd
}

func checkServiceCmd(cmd *types.Cmd, state *types.State) (*types.Cmd, bool) {
	/*
		Checks if a service is valid from a command
	*/
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
	/*
		Get command for service
	*/
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

func sendCommand(cmd *types.Cmd, state *types.State) *terr.Trace {
	/*
		Sends the command results back to the client
	*/
	if cmd.Trace != nil {
		var rvs []interface{}
		rvs = append(rvs, cmd.Trace.Log())
		cmd.ReturnValues = rvs
		cmd.Status = "error"
	} else {
		cmd.Status = "success"
	}
	payload, err := json.Marshal(cmd)
	if err != nil {
		msg := "Unable to marshall json: " + err.Error()
		err := errors.New(msg)
		trace := terr.New(err)
		cmd.Trace = trace
		return trace
	}
	_, err = state.Cli.Http.Publish(state.WsServer.CmdChanOut, payload)
	if err != nil {
		cmd.Trace = terr.New(err)
		return cmd.Trace
	}
	return nil
}
