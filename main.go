package main

import (
	"time"
    "net/http"
    "log"
    "github.com/julienschmidt/httprouter"
    "github.com/acmacalister/skittles"
    "github.com/synw/microb/libmicrob/events"
    "github.com/synw/microb/libmicrob/metadata"
    "github.com/synw/microb/libmicrob/http_handlers"
    "github.com/synw/microb/libmicrob/listeners"
)


func init() {
	go listeners.ListenToIncomingCommands()
}

func main() {
    router := httprouter.New()
    router.GET("/", http_handlers.ServeRequest)
    router.GET("/p/*url", http_handlers.ServeRequest)
    router.GET("/x/*url", http_handlers.ServeApi)
    router.ServeFiles("/static/*filepath", http.Dir("static"))
    router.PanicHandler = http_handlers.Handle500
    server := metadata.GetServer()
    database := metadata.GetMainDatabase()
    loc := server.Host+":"+server.Port
    if (metadata.GetVerbosity() > 0) {
		msg := "Server started on "+loc+" for domain "+skittles.BoldWhite(server.Domain)
		msg = msg+" with "+database.Type
		msg = msg+" ("+database.Host+")"
		events.New("runtime_info", "http_server", msg)
	}
	httpServer := &http.Server{
		Addr: loc,
	    ReadTimeout: 5 * time.Second,
	    WriteTimeout: 10 * time.Second,
	    Handler: router,
	}
    log.Fatal(httpServer.ListenAndServe())
}
