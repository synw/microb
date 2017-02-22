package http_handlers

import (
	"fmt"
	"strings"
	"errors"
	"net/http"
	"encoding/json"
    "html/template"
    "github.com/acmacalister/skittles"
    "github.com/synw/microb/libmicrob/events"
    "github.com/synw/microb/libmicrob/datatypes"
    "github.com/synw/microb/libmicrob/db"
    "github.com/synw/microb/libmicrob/state"
)


var View = template.Must(template.New("view.html").ParseFiles("templates/view.html", "templates/head.html", "templates/header.html", "templates/navbar.html", "templates/footer.html", "templates/routes.js"))
var V404 = template.Must(template.New("404.html").ParseFiles("templates/404.html", "templates/head.html", "templates/header.html", "templates/navbar.html", "templates/footer.html", "templates/routes.js"))
var V500 = template.Must(template.New("500.html").ParseFiles("templates/500.html", "templates/head.html", "templates/header.html", "templates/navbar.html", "templates/footer.html", "templates/routes.js"))


func StartMsg() string {
	return startMsg()
}

func ServeRequest(response http.ResponseWriter, request *http.Request) {
	url := request.URL.Path
	d := make(map[string]interface{})
	d["status_code"] = http.StatusOK
	status := http.StatusOK
    page := &datatypes.Page{Url: url, Title: "", Content: ""}
    response = httpResponseWriter{response, &status}
    renderTemplate(response, page)
}

func ServeApi(response http.ResponseWriter, request *http.Request) {
	url := request.URL.Path
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
	fmt.Fprintf(response, "%s\n", json_bytes)
}

func ReparseTemplates() {
	View = template.Must(template.New("view.html").ParseFiles("templates/view.html", "templates/head.html", "templates/header.html", "templates/navbar.html", "templates/footer.html", "templates/routes.js"))
}

func Handle404(response http.ResponseWriter, request *http.Request) {
	d := make(map[string]interface{})
	d["status_code"] = http.StatusNotFound
	msg := "Not found"
	event := &datatypes.Event{"request_error", "http_server", msg, d}
	events.Handle(event)
	page := &datatypes.Page{Url: "/error/", Title: "Page not found", Content: "<h1>Page not found</h1>"}
	status := http.StatusNotFound
	response = httpResponseWriter{response, &status}
	render404(response, page)
}

func handleAPI404(response http.ResponseWriter, request *http.Request) {
	d := make(map[string]interface{})
	d["status_code"] = http.StatusNotFound
	msg := "Not found"
	event := &datatypes.Event{"request_error", "http_server", msg, d}
	events.Handle(event)
	page := &datatypes.Page{Url: "/error/", Title: "Page not found", Content: "<h1>Page not found</h1>"}
	status := http.StatusNotFound
	response = httpResponseWriter{response, &status}
	json_bytes, _ := json.Marshal(page)
	fmt.Fprintf(response, "%s\n", json_bytes)
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
	for _, route := range(state.Routes) {
		if route == url {
			is_valid = true
			break
		}
		apiroute := "/x/"+route
		if apiroute == url {
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
	database := state.Server.PagesDb
	server := state.Server
    loc := server.Host+":"+server.Port
    if (state.Verbosity > 0) {
		msg = "Server started on "+loc+" for domain "+skittles.BoldWhite(server.Domain)
		msg = msg+" with "+database.Type
		msg = msg+" ("+database.Host+")"
		events.New("runtime_info", "http_server", msg)
	}
	return msg
}

func formatUrl(url string) string {
	if strings.HasPrefix(url, "/") == false {
		url = "/"+url
	}
	return url
}
