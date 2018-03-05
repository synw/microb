package libmicrob

import (
	"github.com/synw/microb/libmicrob/types"
)

func handle(event *types.Event) {
	if Verbose() == true {
		PrintEvent(event)
	}
}
