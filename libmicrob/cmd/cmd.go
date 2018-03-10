package cmd

import (
	//sid "github.com/SKAhack/go-shortid"
	"errors"
	"github.com/synw/microb/libmicrob/types"
	"github.com/synw/terr"
	"time"
)

func New(name string, service *types.Service, from string, args ...interface{}) *types.Cmd {
	// TO FIX
	//id, _ := sid.Generate()
	id := "id"
	date := time.Now()
	status := "pending"
	var tr *terr.Trace
	var rvs []interface{}
	var exec func(*types.Cmd) *types.Cmd
	cmd := &types.Cmd{
		id,
		name,
		date,
		service,
		args,
		from,
		status,
		"",
		tr,
		rvs,
		exec,
	}
	return cmd
}

func ConvertPayload(payload interface{}) (*types.Cmd, bool) {
	pl := payload.(map[string]interface{})
	status := pl["Status"].(string)
	name := pl["Name"].(string)
	//serv := pl["Service"].(string)
	from := pl["From"].(string)
	var args []interface{}
	if pl["Args"] != nil {
		args = pl["Args"].([]interface{})
	}
	cmd := &types.Cmd{
		Name: name,
		//Service: serv,
		From:   from,
		Args:   args,
		Status: status,
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
	if status != "pending" {
		return cmd, false
	}
	return cmd, true
}
