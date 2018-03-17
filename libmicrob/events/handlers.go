package events

import (
	"github.com/synw/microb/libmicrob/msgs"
	"github.com/synw/microb/libmicrob/types"
)

func handle(event *types.Event) {
	msgs.PrintEvent(event)
}
