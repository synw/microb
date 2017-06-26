package cmd

import (
	"github.com/synw/microb/libmicrob/types"
	"github.com/synw/microb/services/httpServer/state/mutate"
	"github.com/synw/terr"
)

func Dispatch(cmd *types.Command) *types.Command {
	com := &types.Command{}
	// TODO: error handling
	if cmd.Name == "start" {
		res := Start(cmd)
		return res
	} else if cmd.Name == "stop" {
		return Stop(cmd)
	}
	return com
}

func Start(cmd *types.Command) *types.Command {
	tr := mutate.StartHttpServer()
	if tr != nil {
		cmd.Trace = tr
		cmd.Status = "error"
		terr.Debug("cmd err", tr)
		return cmd
	}
	var resp []interface{}
	resp = append(resp, "Http server started")
	cmd.Status = "success"
	cmd.ReturnValues = resp
	return cmd
}

func Stop(cmd *types.Command) *types.Command {
	tr := mutate.StopHttpServer()
	if tr != nil {
		cmd.Trace = tr
		cmd.Status = "error"
		return cmd
	}
	var resp []interface{}
	resp = append(resp, "Http server stopped")
	cmd.Status = "success"
	cmd.ReturnValues = resp
	return cmd
}
