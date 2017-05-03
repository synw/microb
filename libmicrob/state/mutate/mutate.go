package mutate

import (
	"github.com/synw/microb/libmicrob/httpServer"
)


func StartHttpServer() {
	go httpServer.Run()
	return
}

func StopHttpServer() {
	_ = httpServer.Stop()
	return
}
