package cli

import (
	"github.com/abiosoft/ishell"
	"github.com/synw/microb-cli/libmicrob/cmd/handler"
	command "github.com/synw/microb/libmicrob/cmd"
	"github.com/synw/terr"
)

func Cmds() *ishell.Cmd {
	command := &ishell.Cmd{
		Name: "http",
		Help: "Http server commands: start, stop",
		Func: func(ctx *ishell.Context) {
			if len(ctx.Args) == 0 {
				err := terr.Err("A parameter is required: ex: http start")
				ctx.Println(err.Error())
				return
			}
			if ctx.Args[0] == "start" {
				cmd := command.New("start", "http", "cli", "")
				cmd, timeout, tr := handler.SendCmd(cmd, ctx)
				if tr != nil {
					tr = terr.Pass("cmd.httpServer.Start", tr)
					msg := tr.Formatc()
					ctx.Println(msg)
				}
				if timeout == true {
					err := terr.Err("Timeout: server is not responding")
					ctx.Println(err.Error())
				}
			} else if ctx.Args[0] == "stop" {
				cmd := command.New("stop", "http", "http", "cli", "")
				cmd, timeout, tr := handler.SendCmd(cmd, ctx)
				if tr != nil {
					tr = terr.Pass("cmd.httpServer.Stop", tr)
					msg := tr.Formatc()
					ctx.Println(msg)
				}
				if timeout == true {
					err := terr.Err("Timeout: server is not responding")
					ctx.Println(err.Error())
				}
			}
			return
		},
	}
	return command
}

/*
func Http() *ishell.Cmd {
	command := &ishell.Cmd{
		Name: "http",
		Help: "Retrieve the state of the http server",
		Func: func(ctx *ishell.Context) {
			cmd := command.New("http", "info", "cli", "")
			cmd, timeout, trace := handler.SendCmd(cmd, ctx)
			if trace != nil {
				trace = terr.Pass("cmd.info.Http", trace)
				msg := trace.Formatc()
				ctx.Println(msg)
				return
			}
			if timeout == true {
				err := terr.Err("Timeout: server is not responding")
				ctx.Println(err.Error())
			}
			return
		},
	}
	return command
}*/
