package info

import (
	"github.com/synw/microb/libmicrob/state"
	"github.com/synw/microb/libmicrob/datatypes"
)


func Http(cmd *datatypes.Command) *datatypes.Command {
	var resp []interface{}
	if (state.HttpServer.Running == true) {
		resp = append(resp, "Http server is up")
	} else {
		resp = append(resp, "Http server is down")
	}
	cmd.Status = "success"
	cmd.ReturnValues = resp
	return cmd
}
