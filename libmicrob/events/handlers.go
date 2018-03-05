package libmicrob

import (
	m "github.com/synw/microb/libmicrob"
	"github.com/synw/microb/libmicrob/types"
)

func handle(event *types.Event) {
	if m.Verbose() == true {
		m.PrintEvent(event)
	}
}
