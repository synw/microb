package httpServer

import (
	"context"
	"encoding/json"
	//"errors"
	"fmt"
	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
	"net/http"
	"strconv"
	//"strings"
	"time"
	//"github.com/goware/cors"
	"github.com/acmacalister/skittles"
	//"github.com/synw/microb/libmicrob/db"
	"github.com/synw/microb/libmicrob/events"
	globalState "github.com/synw/microb/libmicrob/state"
	"github.com/synw/microb/services/httpServer/datatypes"
	"github.com/synw/microb/services/httpServer/state"
	"github.com/synw/terr"
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
	events.Msg("state", "httpServer.run", startMsg())
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
		return tr
	}
	state.HttpServer.Running = false
	events.Msg("state", "httpServer.Stop", stopMsg())
	return nil
}

func handle404(response http.ResponseWriter, request *http.Request) {
	/*msg := "Document not found"
	events.Msg("request_error", "httpServer.handle404", msg)
	fields := make(map[string]interface{})
	fields["Title"] = "Page not found"
	fields["Content"] = "<h1>Page not found</h1>"
	doc := &datatypes.Document{Url: "/error/", Fields: fields}
	status := http.StatusNotFound
	response = httpResponseWriter{response, &status}
	response = headers(response)
	json_bytes, _ := json.Marshal(doc)
	fmt.Fprintf(response, "%s\n", json_bytes)*/
}

func ServeApi(response http.ResponseWriter, request *http.Request) {
	url := request.URL.Path
	if url == "" {
		url = "/"
	}
	doc, err := getDocument(url)
	if err != nil {
		events.Err("error", "httpServer.ServeApi", err)
	}
	if doc == nil {
		if globalState.Debug == true {
			fmt.Println("http.handlers.ServeApi() error: route " + url + " not found from database")
		}
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
	/*index_url := url
	found := false
	// remove url mask
	index_url = strings.Replace(index_url, "/x", "", -1)
	// hit db
	doc, found, tr := db.GetByUrl(index_url)
	if tr != nil {
		tr := terr.Pass("httpServer.getDocument", tr)
		return doc, tr
	}
	if found == false {
		msg := "Document " + url + " not found"
		err := errors.New(msg)
		tr := terr.Add("httpServer.getDocument", err, tr)
		return doc, tr
	}*/
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
	/*database := state.HttpServer.PagesDb
		server := state.HttpServer
	    loc := server.Host+":"+server.Port
		msg = "Server started on "+loc+" for domain "+skittles.BoldWhite(server.Domain)
		msg = msg+" with "+database.Type
		msg = msg+" ("+database.Host+")"
		events.New("runtime_info", "http_server", msg)*/
	return msg
}
