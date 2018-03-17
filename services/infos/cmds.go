package infos

import (
	"github.com/synw/microb/libmicrob/types"
)

func getCmds() map[string]*types.Cmd {
	cmds := make(map[string]*types.Cmd)
	cmds["ping"] = ping()
	cmds["services"] = srv()
	return cmds
}

func ping() *types.Cmd {
	cmd := &types.Cmd{Name: "ping", Exec: runPing}
	return cmd
}

func srv() *types.Cmd {
	cmd := &types.Cmd{Name: "services", Exec: runSrv}
	return cmd
}

func runPing(cmd *types.Cmd, c chan *types.Cmd, state *types.State) {
	var resp []interface{}
	resp = append(resp, "PONG")
	cmd.ReturnValues = resp
	cmd.Status = "success"
	c <- cmd
}

func runSrv(cmd *types.Cmd, c chan *types.Cmd, state *types.State) {
	s := "Available services: "
	for _, v := range state.Services {
		s = s + v.Name + " "
	}
	var rvs []interface{}
	rvs = append(rvs, s)
	cmd.ReturnValues = rvs
	cmd.Status = "success"
	c <- cmd
}
