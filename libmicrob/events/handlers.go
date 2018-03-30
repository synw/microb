package events

import (
	"github.com/synw/microb/libmicrob/msgs"
	"github.com/synw/microb/libmicrob/types"
	"github.com/synw/microb/services/logs"
)

func handle(event *types.Event) {
	msgs.PrintEvent(event)
	logs.Event(event)
}
