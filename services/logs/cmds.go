package logs

import (
	color "github.com/acmacalister/skittles"
	"github.com/synw/microb/libmicrob/types"
	"github.com/vjeantet/jodaTime"
)

func getCmds() map[string]*types.Cmd {
	cmds := make(map[string]*types.Cmd)
	cmds["get"] = get()
	return cmds
}

func get() *types.Cmd {
	cmd := &types.Cmd{Name: "get", Exec: runGet}
	return cmd
}

func runGet(cmd *types.Cmd, c chan *types.Cmd) {
	var resp []interface{}
	logData := getLogs(10)
	msg := "Found logs:\n"
	resp = append(resp, msg)
	for _, log := range logData {
		date := jodaTime.Format("dd/MM/YYYY HH'h'mm:ss", log.CreatedAt)
		msg = date + " [" + log.Level + "] " + color.Blue(log.Service)
		msg = msg + " <" + color.Yellow(log.Class) + "> " + log.Msg
		resp = append(resp, msg)
	}
	cmd.ReturnValues = resp
	cmd.Status = "success"
	c <- cmd
}
