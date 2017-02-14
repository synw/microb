package main

import (
    "fmt"
    "net/http"
    "log"
    "strings"
    "encoding/json"
    "html/template"
    "github.com/julienschmidt/httprouter"
    "github.com/acmacalister/skittles"
    "github.com/synw/microb/events"
    "github.com/synw/microb/datatypes"
    "github.com/synw/microb/db"
    "github.com/synw/microb/metadata"
)


var Routes = db.GetRoutes()
var View = template.Must(template.New("view.html").ParseFiles("templates/view.html", "templates/routes.js"))

func renderTemplate(response http.ResponseWriter, page *datatypes.Page) {
    err := View.Execute(response, page)
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

func getPage(url string) *datatypes.Page {
	hasSlash := strings.HasSuffix(url, "/")
	index_url := url
	found := false
	var data map[string]interface{}
	page := &datatypes.Page{Url:"404", Title:"404", Content:"404 Page not found"}
	if (hasSlash == false) {
		index_url = url+"/"
	}
	// remove url mask
	index_url = strings.Replace(index_url,"/x","",-1)
	// hit db
	data, found = db.GetFromUrl(index_url)
	if (found == false) {
		return page
	}
	fields := data["fields"].(map[string]interface{})
	content := fields["content"].(string)
	title := fields["title"].(string)
	page = &datatypes.Page{Url: url, Title: title, Content: content}
	return page
}

func serveRequest(response http.ResponseWriter, request *http.Request, ps httprouter.Params) {
	url := ps.ByName("url")
	if url == "" {url = "/"}
	if isValidRoute(url) == false {
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

type httpResponseWriter struct {
	http.ResponseWriter
	status *int
}

func serveApi(response http.ResponseWriter, request *http.Request, ps httprouter.Params) {
	url := ps.ByName("url")
	page := getPage(url)
	if isValidRoute(url) == false {
		handle404(response, request, url, true)
		return
	}
	if (page.Url == "404") {
    	handle404(response, request, url, true)
    	return
    }
	msg := "API "+url
	d := make(map[string]interface{})
	d["status_code"] = http.StatusOK
	event := &datatypes.Event{"request", "http_server", msg, d}
    events.Handle(event)
    status := http.StatusOK
	json_bytes, _ := json.Marshal(page)
	response = httpResponseWriter{response, &status}
	fmt.Fprintf(response, "%s\n", json_bytes)
}

func handle500(response http.ResponseWriter, request *http.Request, params interface{}) {
	msg := "Error 500" 
	event := events.NewEvent("error", "http_server", msg)
	events.Handle(event)
	page := &datatypes.Page{Url: "/error/", Title: "Technical error", Content: "<h1>Technical error</h1>"}
	status := http.StatusInternalServerError
	json_bytes, _ := json.Marshal(page)
	response = httpResponseWriter{response, &status}
	fmt.Fprintf(response, "%s\n", json_bytes)
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
		renderTemplate(response, page)
	}
}

func main() {
    router := httprouter.New()
    router.GET("/", serveRequest)
    router.GET("/p/*url", serveRequest)
    router.GET("/x/*url", serveApi)
    router.ServeFiles("/static/*filepath", http.Dir("static"))
    router.PanicHandler = handle500
    server := metadata.GetServer()
    database := metadata.GetMainDatabase()
    loc := server.Host+":"+server.Port
    if (metadata.GetVerbosity() > 0) {
		msg := "Server started on "+loc+" for domain "+skittles.BoldWhite(server.Domain)
		msg = msg+" with "+database.Type
		msg = msg+" ("+database.Host+")"
		event := events.NewEvent("runtime_info", "http_server", msg)
		events.Handle(event)
	}
	
    log.Fatal(http.ListenAndServe(loc, router))
}
