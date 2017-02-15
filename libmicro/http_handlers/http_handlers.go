package http_handlers

import (
	"fmt"
	"strings"
	"net/http"
	"encoding/json"
    "html/template"
    "github.com/julienschmidt/httprouter"
    "github.com/synw/microb/libmicro/events"
    "github.com/synw/microb/libmicro/datatypes"
    "github.com/synw/microb/libmicro/db"
)


var Routes = db.GetRoutes()
var View = template.Must(template.New("view.html").ParseFiles("templates/view.html", "templates/routes.js"))
var V404 = template.Must(template.New("404.html").ParseFiles("templates/404.html", "templates/routes.js"))
var V500 = template.Must(template.New("500.html").ParseFiles("templates/500.html", "templates/routes.js"))

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

func ServeRequest(response http.ResponseWriter, request *http.Request, ps httprouter.Params) {
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

func ServeApi(response http.ResponseWriter, request *http.Request, ps httprouter.Params) {
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

func Handle500(response http.ResponseWriter, request *http.Request, params interface{}) {
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
