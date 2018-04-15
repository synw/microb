package logs

import (
	color "github.com/acmacalister/skittles"
	"github.com/synw/microb/libmicrob/types"
	"github.com/vjeantet/jodaTime"
)

func getCmds() map[string]*types.Cmd {
	cmds := make(map[string]*types.Cmd)
	cmds["get"] = get()
	cmds["errors"] = errs()
	cmds["warnings"] = warns()
	cmds["state"] = state()
	cmds["cmds"] = eXcommands()
	return cmds
}

func get() *types.Cmd {
	cmd := &types.Cmd{Name: "get", Exec: runGet}
	return cmd
}

func errs() *types.Cmd {
	cmd := &types.Cmd{Name: "errors", Exec: runErrs}
	return cmd
}

func warns() *types.Cmd {
	cmd := &types.Cmd{Name: "warnings", Exec: runWarns}
	return cmd
}

func state() *types.Cmd {
	cmd := &types.Cmd{Name: "state", Exec: runState}
	return cmd
}

func eXcommands() *types.Cmd {
	cmd := &types.Cmd{Name: "cmds", Exec: runDbCommands}
	return cmd
}

func runState(cmd *types.Cmd, c chan *types.Cmd, args ...interface{}) {
	var resp []interface{}
	logData := getState(10)
	msg := "State log:"
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

func runDbCommands(cmd *types.Cmd, c chan *types.Cmd, args ...interface{}) {
	var resp []interface{}
	logData := getCommandsFromDb(10)
	msg := "Commands log:"
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

func runWarns(cmd *types.Cmd, c chan *types.Cmd, args ...interface{}) {
	var resp []interface{}
	logData := getWarns(10)
	msg := "Found warnings:"
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

func runErrs(cmd *types.Cmd, c chan *types.Cmd, args ...interface{}) {
	var resp []interface{}
	logData := getErrs(10)
	msg := "Found errors:"
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

func runGet(cmd *types.Cmd, c chan *types.Cmd, args ...interface{}) {
	var resp []interface{}
	logData := getLogs(10)
	msg := "Found logs:"
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
