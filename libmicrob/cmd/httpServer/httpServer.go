package httpServer

import (
	"github.com/synw/microb/libmicrob/datatypes"
	"github.com/synw/microb/libmicrob/state/mutate"
)


func Start(cmd *datatypes.Command) *datatypes.Command {
	mutate.StartHttpServer()
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
