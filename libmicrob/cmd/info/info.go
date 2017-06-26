package base

import (
	"github.com/synw/microb/libmicrob/types"
)

func Dispatch(cmd *types.Command) *types.Command {
	com := &types.Command{}
	// TODO: error handling
	if cmd.Name == "ping" {
		return Ping(cmd)
	}
	return com
}

func Ping(cmd *types.Command) *types.Command {
	var resp []interface{}
	resp = append(resp, "PONG")
	cmd.ReturnValues = resp
	cmd.Status = "success"
	return cmd
}
