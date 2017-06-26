package info

import (
	"github.com/synw/microb/libmicrob/types"
	"github.com/synw/microb/services/httpServer/state"
)

func Http(cmd *types.Command) *types.Command {
	var resp []interface{}
	if state.HttpServer.Running == true {
		resp = append(resp, "Http server is up")
	} else {
		resp = append(resp, "Http server is down")
	}
	cmd.Status = "success"
	cmd.ReturnValues = resp
	return cmd
}
