package httpServer

import (
	"github.com/synw/microb/libmicrob/datatypes"
	"github.com/synw/microb/libmicrob/state/mutate"
	"github.com/synw/microb/libmicrob/events"
)


func Start(cmd *datatypes.Command) *datatypes.Command {
	tr := mutate.StartHttpServer()
	if tr != nil {
		events.Err("error", "cmd.httpServer", tr)
		return cmd
	}
	var resp []interface{}
	resp = append(resp, "Http server started")
	cmd.Status = "success"
	cmd.ReturnValues = resp
	return cmd
}

func Stop(cmd *datatypes.Command) *datatypes.Command {
	mutate.StopHttpServer()
	var resp []interface{}
	resp = append(resp, "Http server stopped")
	cmd.Status = "success"
	cmd.ReturnValues = resp
	return cmd
}
