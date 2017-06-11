package httpServer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/acmacalister/skittles"
	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
	"github.com/synw/microb/libmicrob/events"
	"github.com/synw/microb/services/httpServer/datatypes"
	"github.com/synw/microb/services/httpServer/state"
	"github.com/synw/terr"
	"net/http"
	"strconv"
	"time"
)

type httpResponseWriter struct {
	http.ResponseWriter
	status *int
}

func InitHttpServer(serve bool) {
	// routing
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.StripSlashes)
	// main route
	r.Route("/", func(r chi.Router) {
		r.Get("/*", ServeApi)
	})
	// init
	loc := state.HttpServer.Host + ":" + strconv.Itoa(state.HttpServer.Port)
	httpServer := &http.Server{
		Addr:         loc,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      r,
	}
	state.HttpServer.Instance = httpServer
	if serve == true {
		Run()
	}
}

func Run() {
	events.New("state", "http", "httpServer.Run", startMsg(), nil)
	state.HttpServer.Running = true
	state.HttpServer.Instance.ListenAndServe()
}

func Stop() *terr.Trace {
	d := time.Now().Add(5 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), d)
	defer cancel()
	srv := state.HttpServer.Instance
	err := srv.Shutdown(ctx)
	if err != nil {
		tr := terr.New("httpServer.Stop", err)
		events.New("error", "http", "httpServer.Stop", stopMsg(), nil)
		return tr
	}
	state.HttpServer.Running = false
	events.New("state", "http", "httpServer.Stop", stopMsg(), nil)
	return nil
}

func handle404(response http.ResponseWriter, request *http.Request) {
	fmt.Println("404")
}

func ServeApi(response http.ResponseWriter, request *http.Request) {
	url := request.URL.Path
	if url == "" {
		url = "/"
	}
	doc, tr := getDocument(url)
	if tr != nil {
		msg := "Unable to get document from database"
		events.Err("http", "httpServer.ServeApi", msg, tr.ToErr())
	}
	if doc == nil {
		msg := "http.handlers.ServeApi() error: route " + url + " not found from database"
		err := errors.New(msg)
		events.Err("http", "httpServer.ServeApi", msg, err)
		handle404(response, request)
		return
	}
	json_bytes, _ := json.Marshal(doc)
	response = headers(response)
	fmt.Fprintf(response, "%s\n", json_bytes)
}

func headers(response http.ResponseWriter) http.ResponseWriter {
	response.Header().Set("Content-Type", "application/json")
	response.Header().Set("Access-Control-Allow-Origin", state.HttpServer.Cors)
	return response
}

func getDocument(url string) (*datatypes.Document, *terr.Trace) {
	doc := &datatypes.Document{}
	return doc, nil
}

func stopMsg() string {
	msg := "Http server stopped"
	return msg
}

func startMsg() string {
	var msg string
	loc := state.HttpServer.Host + ":" + strconv.Itoa(state.HttpServer.Port)
	msg = "Http server started at " + loc + " for domain " + skittles.BoldWhite(state.HttpServer.Domain)
	return msg
}
