package info

import (
	"github.com/synw/microb/libmicrob/types"
)

func Dispatch(cmd *types.Cmd) *types.Cmd {
	com := &types.Cmd{}
	if cmd.Name == "ping" {
		return Ping(cmd)
	}
	return com
}

func Ping(cmd *types.Cmd) *types.Cmd {
	var resp []interface{}
	resp = append(resp, "PONG")
	cmd.ReturnValues = resp
	cmd.Status = "success"
	return cmd
}
