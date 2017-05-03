package httpServer

import (
	"net/http"
	"fmt"
	"strconv"
	"time"
	"context"
	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
	"github.com/acmacalister/skittles"
	"github.com/synw/microb/libmicrob/datatypes"
	"github.com/synw/microb/libmicrob/state"
	//"github.com/synw/microb/libmicrob/events"
)


func InitHttpServer() {
	// routing
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.StripSlashes)
	r.Get("/x/", ServeApi)
	// init
	loc := state.Server.HttpHost+":"+strconv.Itoa(state.Server.HttpPort)
	httpServer := &http.Server{
		Addr: loc,
	    ReadTimeout: 5 * time.Second,
	    WriteTimeout: 10 * time.Second,
	    Handler: r,
	}
	c := make(chan struct{})
	state.HttpServer = datatypes.HttpServer{state.Server, httpServer, false, c}
	return 
}

func Run() {
	fmt.Println(startMsg())
	state.HttpServer.Instance.ListenAndServe()
	state.HttpServer.Runing = true
}

func Stop() error {
	d := time.Now().Add(5 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), d)
	defer cancel()
	srv := state.HttpServer.Instance
	err := srv.Shutdown(ctx)
	if err != nil {
		return err
	}
	state.HttpServer.Runing = false
	fmt.Println(stopMsg())
	return nil
}

func ServeApi(response http.ResponseWriter, request *http.Request) {
	/*url := request.URL.Path
	if url == "/x" {
		url = "/"
	} 
	doc, err := getDocument(url)
	if err != nil {
		events.Error("http_handlers.ServeApi()", err)
	}
	if (doc == nil) {
		if state.Debug == true {
			fmt.Println("http.handlers.ServeApi() error: route "+url+" not found from database")
		}
    	handleAPI404(response, request)
    	return
    }
	json_bytes, _ := json.Marshal(doc)
	fmt.Fprintf(response, "%s\n", json_bytes)*/
	fmt.Println("SERVE", request.URL.Path)
}
/*
func getDocument(url string) (*datatypes.Page, error) {
	index_url := url
	found := false
	// remove url mask
	index_url = strings.Replace(index_url,"/x","",-1)
	// hit db
	doc, found, err := db.GetFromUrl(index_url)
	if err != nil {
		events.Error("http_handlers.getDocument", err)
		return doc, err
	}
	if (found == false) {
		msg := "Document "+url+" not found"
		err = errors.New(msg)
		events.Error("http_handlers.getDocument", err)
	}
	return doc, nil
}
*/

func stopMsg() string {
	// welcome msg
	msg := "Http server stopped"
	/*database := state.HttpServer.PagesDb
	server := state.HttpServer
    loc := server.Host+":"+server.Port
	msg = "Server started on "+loc+" for domain "+skittles.BoldWhite(server.Domain)
	msg = msg+" with "+database.Type
	msg = msg+" ("+database.Host+")"
	events.New("runtime_info", "http_server", msg)*/
	return msg
}

func startMsg() string {
	// welcome msg
	var msg string
    loc := state.Server.HttpHost+":"+strconv.Itoa(state.Server.HttpPort)
	msg = "Server started on "+loc+" for domain "+skittles.BoldWhite(state.Server.Domain)
	/*database := state.HttpServer.PagesDb
	server := state.HttpServer
    loc := server.Host+":"+server.Port
	msg = "Server started on "+loc+" for domain "+skittles.BoldWhite(server.Domain)
	msg = msg+" with "+database.Type
	msg = msg+" ("+database.Host+")"
	events.New("runtime_info", "http_server", msg)*/
	return msg
}
