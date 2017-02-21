package main

import (
	"time"
    "net/http"
    "github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
    "github.com/synw/microb/libmicrob/http_handlers"
    "github.com/synw/microb/libmicrob/listeners"
    "github.com/synw/microb/libmicrob/state"
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
	r.Route("/x", func(r chi.Router) {
		r.Get("/", http_handlers.ServeApi)
	    r.Route("/:path", func(r chi.Router) {
	      	r.Get("/", http_handlers.ServeApi)
		})
	})
	r.Get("/:url", http_handlers.ServeRequest)
	r.Get("/", http_handlers.ServeRequest)
	// http server
	http_handlers.StartMsg()
	loc := state.Server.Host+":"+state.Server.Port
    if state.Server.PagesDb != nil {
		httpServer := &http.Server{
			Addr: loc,
		    ReadTimeout: 5 * time.Second,
		    WriteTimeout: 10 * time.Second,
		    Handler: r,
		}
	    httpServer.ListenAndServe()
	    state.Server.Runing = true
	} else {
		// sit
		run := make(chan bool)
		<- run
	}
}
