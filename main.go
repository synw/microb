package main

import (
	"time"
    "net/http"
    "github.com/pressly/chi"
	//"github.com/pressly/chi/docgen"
	"github.com/pressly/chi/middleware"
	//"github.com/pressly/chi/render"
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
	// routing
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Get("/x/", http_handlers.ServeApi)
	//r.Get("/x/:path", http_handlers.ServeApi)
	
	r.Route("/x", func(r chi.Router) {
		r.Get("/", http_handlers.ServeApi)
	    r.Route("/:path", func(r chi.Router) {
	      	r.Get("/", http_handlers.ServeApi)
		})
	  })
	
	r.Get("/:url", http_handlers.ServeRequest)
	r.Get("/", http_handlers.ServeRequest)
    // welcome msg
    server := metadata.GetServer()
    database := metadata.GetMainDatabase()
    loc := server.Host+":"+server.Port
    if (metadata.GetVerbosity() > 0) {
		msg := "Server started on "+loc+" for domain "+skittles.BoldWhite(server.Domain)
		msg = msg+" with "+database.Type
		msg = msg+" ("+database.Host+")"
		events.New("runtime_info", "http_server", msg)
	}
	// http server
	httpServer := &http.Server{
		Addr: loc,
	    ReadTimeout: 5 * time.Second,
	    WriteTimeout: 10 * time.Second,
	    Handler: r,
	}
    httpServer.ListenAndServe()
}
