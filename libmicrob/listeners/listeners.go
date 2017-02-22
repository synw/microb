package listeners

import (
	"sync"
	"github.com/synw/microb/libmicrob/listeners/websockets"
	"github.com/synw/microb/libmicrob/events"
	"github.com/synw/microb/libmicrob/state"
)


func ListenToIncomingCommands(wg *sync.WaitGroup) {
	channel_name := "$"+state.Server.Domain+"_commands"
	c_err := make(chan error)
	c_ok := make(chan bool)
	var is_ok bool
	go websockets.ListenToIncomingCommands(channel_name, c_ok, c_err)
	go func() {
		select {
			case is_ok = <- c_ok:
				if is_ok == false {
					state.ListenWs = false
					defer wg.Done()
					
				} else {
					state.ListenWs = true
					defer wg.Done()
					msg := "Listening to websockets channel "+channel_name
					events.New("info", "listeners.ListenToIncomingCommands", msg)
				}
				close(c_ok)
		}
	}()
	var err error
	go func() {
		select {
			case err = <- c_err:
				msg := "Unable to connect to websockets server: "+err.Error()
				events.ErrMsg("listeners.ListenToIncomingCommands", msg)
				close(c_err)
		}
	}()
}
