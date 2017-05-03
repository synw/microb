package httpServer

import (
	"github.com/synw/microb/libmicrob/datatypes"
	"github.com/synw/microb/libmicrob/state/mutate"
	"github.com/synw/microb/libmicrob/events"
)


func Start(cmd *datatypes.Command) *datatypes.Command {
	mutate.StartHttpServer()
	var resp []interface{}
	resp = append(resp, "Http server started")
	cmd.Status = "success"
	cmd.ReturnValues = resp
	events.Msg("state", "cmd.httpServer.Start", "Http server started")
	return cmd
}

func Stop(cmd *datatypes.Command) *datatypes.Command {
	mutate.StopHttpServer()
	var resp []interface{}
	resp = append(resp, "Http server stopped")
	cmd.Status = "success"
	cmd.ReturnValues = resp
	events.Msg("state", "cmd.httpServer.Stop", "Http server stoped")
	return cmd
}
