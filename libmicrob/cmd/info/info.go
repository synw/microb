package base

import (
	"github.com/synw/microb/libmicrob/datatypes"
)

func Dispatch(cmd *datatypes.Command) *datatypes.Command {
	com := &datatypes.Command{}
	// TODO: error handling
	if cmd.Name == "ping" {
		return Ping(cmd)
	}
	return com
}

func Ping(cmd *datatypes.Command) *datatypes.Command {
	var resp []interface{}
	resp = append(resp, "PONG")
	cmd.ReturnValues = resp
	cmd.Status = "success"
	return cmd
}
