package http_handlers

import (
	"fmt"
	"strings"
	"errors"
	"net/http"
	"encoding/json"
    "html/template"
    "github.com/pressly/chi"
    "github.com/acmacalister/skittles"
    "github.com/synw/microb/libmicrob/events"
    "github.com/synw/microb/libmicrob/datatypes"
    "github.com/synw/microb/libmicrob/db"
    "github.com/synw/microb/libmicrob/state"
)


var Routes, _ = db.GetRoutes()
var View = template.Must(template.New("view.html").ParseFiles("templates/view.html", "templates/head.html", "templates/header.html", "templates/navbar.html", "templates/footer.html", "templates/routes.js"))
var V404 = template.Must(template.New("404.html").ParseFiles("templates/404.html", "templates/head.html", "templates/header.html", "templates/navbar.html", "templates/footer.html", "templates/routes.js"))
var V500 = template.Must(template.New("500.html").ParseFiles("templates/500.html", "templates/head.html", "templates/header.html", "templates/navbar.html", "templates/footer.html", "templates/routes.js"))


func StartMsg() string {
	return startMsg()
}

func ServeRequest(response http.ResponseWriter, request *http.Request) {
	url := formatUrl(chi.URLParam(request, "url"))
	if isValidRoute(url) == false {
		//fmt.Println("invalid route", url)
		handle404(response, request, url, false)
		return
	}
	msg := "PAGE "+url
	d := make(map[string]interface{})
	d["status_code"] = http.StatusOK
	status := http.StatusOK
	event := &datatypes.Event{"request", "http_server", msg, d}
    events.Handle(event)
    page := &datatypes.Page{Url: url, Title: "", Content: ""}
    response = httpResponseWriter{response, &status}
    renderTemplate(response, page)
}

func ServeApi(response http.ResponseWriter, request *http.Request) {
	url := formatUrl(chi.URLParam(request, "path"))
	doc, err := getDocument(url)
	if isValidRoute(url) == false {
		if state.Debug == true {
			msg := "Invalid route "+url
			msg = msg+doc.Format()
			err = errors.New(msg)
			events.Error("http.handlers.ServeApi", err)
		}
		handle404(response, request, url, true)
		return
	}
	if err != nil {
		events.Error("http_handlers.ServeApi()", err)
	}
	if (doc == nil) {
		if state.Debug == true {
			fmt.Println("http.handlers.ServeApi() error: route "+url+" not found from database")
		}
    	handle404(response, request, url, true)
    	return
    }
	/*
	msg := "API "+url
	d := make(map[string]interface{})
	d["status_code"] = http.StatusOK
	event := &datatypes.Event{"request", "http_server", msg, d}
    events.Handle(event)
    status := http.StatusOK
	json_bytes, _ := json.Marshal(page)
	response = httpResponseWriter{response, &status}
	*/
	json_bytes, _ := json.Marshal(doc)
	fmt.Fprintf(response, "%s\n", json_bytes)
}

func ReparseTemplates() {
	View = template.Must(template.New("view.html").ParseFiles("templates/view.html", "templates/head.html", "templates/header.html", "templates/navbar.html", "templates/footer.html", "templates/routes.js"))
}

func renderTemplate(response http.ResponseWriter, page *datatypes.Page) {
    err := View.Execute(response, page)
    if err != nil {
        http.Error(response, err.Error(), http.StatusInternalServerError)
    }
}

func render404(response http.ResponseWriter, page *datatypes.Page) {
	err := V404.Execute(response, page)
    if err != nil {
        http.Error(response, err.Error(), http.StatusInternalServerError)
    }
}
    
func render500(response http.ResponseWriter, page *datatypes.Page) {
	err := V500.Execute(response, page)
    if err != nil {
        http.Error(response, err.Error(), http.StatusInternalServerError)
    }
}

func isValidRoute(url string) bool {
	is_valid := false
	for _, route := range(Routes) {
		if route == url {
			is_valid = true
			break
		}
	}
	return is_valid
}

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

type httpResponseWriter struct {
	http.ResponseWriter
	status *int
}

func handle404(response http.ResponseWriter, request *http.Request, url string, api bool) {
	d := make(map[string]interface{})
	d["status_code"] = http.StatusNotFound
	var el string
	if (api == true) {
		el = "API"
	} else {
		el = "PAGE"
	}
	msg := el+" "+url+" not found"
	event := &datatypes.Event{"request_error", "http_server", msg, d}
	events.Handle(event)
	page := &datatypes.Page{Url: "/error/", Title: "Page not found", Content: "<h1>Page not found</h1>"}
	status := http.StatusNotFound
	response = httpResponseWriter{response, &status}
	if api == true {
		json_bytes, _ := json.Marshal(page)
		fmt.Fprintf(response, "%s\n", json_bytes)
	} else {
		render404(response, page)
	}
}

/*
func handle500(response http.ResponseWriter, request *http.Request, params interface{}) {
	msg := "Error 500"
	d := make(map[string]interface{})
	d["status_code"] = http.StatusInternalServerError
	event := &datatypes.Event{"request_error", "http_server", msg, d}
	events.Handle(event)
	page := &datatypes.Page{Url: "/error/", Title: "", Content: ""}
	status := http.StatusInternalServerError
	response = httpResponseWriter{response, &status}
	render500(response, page)
}
*/

func startMsg() string {
	// welcome msg
	var msg string
	if state.Server.PagesDb != nil {
		database := state.Server.PagesDb
		server := state.Server
	    loc := server.Host+":"+server.Port
	    if (state.Verbosity > 0) {
			msg = "Server started on "+loc+" for domain "+skittles.BoldWhite(server.Domain)
			msg = msg+" with "+database.Type
			msg = msg+" ("+database.Host+")"
			events.New("runtime_info", "http_server", msg)
		}
	} else {
		events.State("server", "No pages database configured. Http server is not runing")	
	}
	return msg
}

func formatUrl(url string) string {
	if strings.HasPrefix(url, "/") == false {
		url = "/"+url
	}
	return url
}
