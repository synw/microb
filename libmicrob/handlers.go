package libmicrob

import (
	"github.com/synw/microb/libmicrob/types"
)

func handle(event *types.Event) {
	if Verb.State.Current() != "zero" {
		Print(event.Msg, event.Class)
	}
}
