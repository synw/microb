package infos

import (
	"github.com/synw/microb/libmicrob/msgs"
	"github.com/synw/microb/libmicrob/types"
)

func getCmds() map[string]*types.Cmd {
	cmds := make(map[string]*types.Cmd)
	cmds["ping"] = ping()
	cmds["services"] = srv()
	cmds["cmds"] = cmds_()
	return cmds
}

func initService(dev bool, start bool) error {
	return nil
}

func ping() *types.Cmd {
	cmd := &types.Cmd{Name: "ping", Exec: runPing, NoLog: true}
	return cmd
}

func srv() *types.Cmd {
	cmd := &types.Cmd{Name: "services", Exec: runSrv}
	return cmd
}

func cmds_() *types.Cmd {
	cmd := &types.Cmd{Name: "cmds", Exec: runCmds}
	return cmd
}

func runCmds(cmd *types.Cmd, c chan *types.Cmd, args ...interface{}) {
	state := args[0].(*types.State)
	var cmds []string
	cmds = append(cmds, "Found commands:")
	for _, service := range state.Services {
		cmds = append(cmds, "# "+msgs.Bold(service.Name)+" service:")
		for _, cmd := range service.Cmds {
			cmds = append(cmds, "  - "+cmd.Name)
		}
	}
	var rvs []interface{}
	for _, cm := range cmds {
		rvs = append(rvs, cm)
	}
	cmd.ReturnValues = rvs
	cmd.Status = "success"
	c <- cmd
}

func runPing(cmd *types.Cmd, c chan *types.Cmd, args ...interface{}) {
	var resp []interface{}
	resp = append(resp, "PONG")
	cmd.ReturnValues = resp
	cmd.Status = "success"
	c <- cmd
}

func runSrv(cmd *types.Cmd, c chan *types.Cmd, args ...interface{}) {
	state := args[0].(*types.State)
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
