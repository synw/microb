package events

import (
	"github.com/SKAhack/go-shortid"
	"github.com/synw/microb/libmicrob/msgs"
	"github.com/synw/microb/libmicrob/types"
	"github.com/synw/terr"
	"time"
)

var g = shortid.Generator()

func New(msg string, args ...map[string]interface{}) *types.Event {
	class, service, cmd, trace, data := getEventArgs(args...)
	date := time.Now()
	id := g.Generate()
	event := &types.Event{id, class, date, msg, service, cmd, trace, data}
	msgs.Debug(event)
	//handle(event)
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
			msgs.Debug(arg)
			for k, v := range arg {
				if k == "class" {
					eclass = v.(string)
				}
			}
		}
	}
	return eclass, eservice, ecmd, etrace, edata
}
