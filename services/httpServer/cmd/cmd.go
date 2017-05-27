package cmd

import (
	"github.com/synw/microb/libmicrob/datatypes"
	"github.com/synw/microb/services/httpServer/state/mutate"
)

func Start(cmd *datatypes.Command) *datatypes.Command {
	tr := mutate.StartHttpServer()
	if tr != nil {
		cmd.Trace = tr
		cmd.Status = "error"
		return cmd
	}
	var resp []interface{}
	resp = append(resp, "Http server started")
	cmd.Status = "success"
	cmd.ReturnValues = resp
	return cmd
}

func Stop(cmd *datatypes.Command) *datatypes.Command {
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
