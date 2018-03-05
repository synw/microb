package libmicrob

import (
	"github.com/SKAhack/go-shortid"
	"github.com/synw/microb/libmicrob/types"
	"github.com/synw/terr"
	"time"
)

var g = shortid.Generator()

func state(mesg string) *types.Event {
	args := make(map[string]interface{})
	args["class"] = "state"
	event := new_(mesg, args)
	return event
}

func new_(msg string, args ...map[string]interface{}) *types.Event {
	class, service, cmd, trace, data := getEventArgs(args...)
	date := time.Now()
	id := g.Generate()
	event := &types.Event{id, class, date, msg, service, cmd, trace, data}
	handle(event)
	return event
}

func getEventArgs(args ...map[string]interface{}) (class string, service *types.Service, cmd *types.Cmd, trace *terr.Trace, data map[string]interface{}) {
	eclass := "default"
	eservice := &types.Service{}
	ecmd := &types.Cmd{}
	etrace := &terr.Trace{}
	edata := make(map[string]interface{})
	if len(args) > 0 {
		for _, arg := range args {
			for k, v := range arg {
				if k == "class" {
					eclass = v.(string)
				} else if k == "service" {
					eservice = v.(*types.Service)
				} else if k == "command" {
					ecmd = v.(*types.Cmd)
				} else if k == "trace" {
					etrace = v.(*terr.Trace)
				} else if k == "data" {
					edata = v.(map[string]interface{})
				}
			}
		}
	}
	return eclass, eservice, ecmd, etrace, edata
}
