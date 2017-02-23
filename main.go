package main

import (
	"fmt"
	"time"
    "net/http"
    "flag"
    "sync"
    "github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
    "github.com/synw/microb/libmicrob/http_handlers"
    "github.com/synw/microb/libmicrob/listeners"
    "github.com/synw/microb/libmicrob/state"
    "github.com/synw/microb/libmicrob/state/mutate"
    "github.com/synw/microb/libmicrob/events"
    "github.com/synw/microb/libmicrob/db"
)


var dev_mode = flag.Bool("d", false, "Dev mode")

func init() {
	flag.Parse()
	err := state.InitState(*dev_mode)
	if err != nil {
		events.Error("init", err)
	}
	err = db.InitDb()
	if err != nil {
		events.Err("init", "Http server can not run: no pages database", err)
	}
	if state.Server.PagesDb.Running == true {
		if state.Debug == true {
			events.Debug("Getting routes from db")
		}
		routes, err := db.GetRoutes()
		if err != nil {
			events.Err("init", "Can't get the routes out of the database", err)
		} else {
			state.SetRoutes(routes)
		}
	}
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go listeners.ListenToIncomingCommands(&wg)
	// routing
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.StripSlashes)
	r.Get("/x/", http_handlers.ServeApi)
	r.Route("/x", func(r chi.Router) {
		r.Get("/", http_handlers.ServeApi)
		for _, route := range(state.Routes) {
		    r.Route(route, func(r chi.Router) {
		      	r.Get("/", http_handlers.ServeApi)
			})
		}
	})
	for _, route := range(state.Routes) {
		r.Get(route, http_handlers.ServeRequest)
	}
	r.Get("/", http_handlers.ServeRequest)
	r.NotFound(http_handlers.Handle404)
	// http server
	loc := state.Server.Host+":"+state.Server.Port
    /*if state.DbIsOk != true {
    	events.State("main", "Http server is not running. Please check your configuration")
    }*/
	httpServer := &http.Server{
		Addr: loc,
	    ReadTimeout: 5 * time.Second,
	    WriteTimeout: 10 * time.Second,
	    Handler: r,
	}
	state.Server.RunningServer = httpServer
	// wait for all init checks to proceed (are we connected to a websockets server?)
	wg.Wait()
	if state.ListenWs == false {
		msg := "The server is not listening for commands. Please configure your websockets server"
		events.State("runtime", msg)
	}
	if (state.Server.PagesDb.Name != "" && state.ListenWs == false) {
		fmt.Println("Nothing to do: going to sleep")
		return
	}
	_ = mutate.StartServer()
	// sit and listen
	run := make(chan bool)
	<- run
}
