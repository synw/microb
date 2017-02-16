package listeners

import (
	"github.com/synw/microb/libmicrob/metadata"
	"github.com/synw/microb/libmicrob/listeners/websockets"
	"github.com/synw/microb/libmicrob/events"
)



/*
func changefeedsBrocker() bool {
	for _, b := range(metadata.GetCommandsBrockers()) {
		if b == "changefeeds" {
			return true
		}
	}
	return false
}
*/
func ListenToIncomingCommands() {
	if metadata.IsWebsocketsBrocker() == true {
		channel_name := "$"+metadata.GetConfString("domain")+"_commands"
		msg := "Listening to websockets channel "+channel_name
		events.New("info", "runtime", msg)
		websockets.ListenToIncomingCommands(channel_name)
	}
}


