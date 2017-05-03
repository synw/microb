package base

import (
	"github.com/synw/microb/libmicrob/datatypes"
)


func Ping(cmd *datatypes.Command) *datatypes.Command {
	var resp []interface{}
	resp = append(resp, "PONG")
	cmd.ReturnValues = resp
	cmd.Status = "success"
	return cmd
}
