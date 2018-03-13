package infos

import (
	m "github.com/synw/microb/libmicrob"
	"github.com/synw/microb/libmicrob/types"
)

func Dispatch(cmd *types.Cmd) *types.Cmd {
	com := &types.Cmd{}
	if cmd.Name == "ping" {
		return ping(cmd)
	} else if cmd.Name == "services" {
		return srv(cmd)
	}
	return com
}

func ping(cmd *types.Cmd) *types.Cmd {
	var resp []interface{}
	resp = append(resp, "PONG")
	cmd.ReturnValues = resp
	cmd.Status = "success"
	return cmd
}

func srv(cmd *types.Cmd) *types.Cmd {
	s := "Available services: "
	for _, v := range m.Services {
		s = s + v.Name + " "
	}
	var rvs []interface{}
	rvs = append(rvs, s)
	cmd.ReturnValues = rvs
	cmd.Status = "success"
	return cmd
}
